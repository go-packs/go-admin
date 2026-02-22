package admin

import (
	"embed"
	"fmt"
	"github.com/go-packs/go-admin/models"
	"github.com/go-packs/go-admin/resource"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

type PageData struct {
	SiteTitle        string
	Resources        map[string]*resource.Resource
	GroupedResources map[string][]*resource.Resource
	GroupedPages     map[string][]*Page
	CurrentResource  *resource.Resource
	Fields           []resource.Field
	Data             []map[string]interface{}
	Item             map[string]interface{}
	Filters          map[string]string
	User             *models.AdminUser
	Stats            []Stat
	Error            string
	Flash            string
	CSS              template.CSS
	Page, PerPage    int
	TotalPages       int
	TotalCount       int64
	HasPrev, HasNext bool
	PrevPage, NextPage int
	Scopes           []resource.Scope
	CurrentScope     string
	Associations     map[string]AssociationData
	ChartData        []ChartWidget
	SortField        string
	SortOrder        string
	RenderedSidebars map[string]template.HTML
}

type ChartWidget struct {
	ID, Label, Type string
	Labels          []string
	Values          []float64
}

type AssociationData struct {
	Resource *resource.Resource
	Fields   []resource.Field
	Items    []map[string]interface{}
	Options  []map[string]interface{}
}

type Stat struct {
	Label string
	Value int64
}

func (reg *Registry) setFlash(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{Name: "admin_flash", Value: message, Path: "/admin", HttpOnly: true})
}

func (reg *Registry) getFlash(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("admin_flash")
	if err != nil { return "" }
	http.SetCookie(w, &http.Cookie{Name: "admin_flash", Value: "", Path: "/admin", MaxAge: -1})
	return cookie.Value
}

// ServeHTTP implements the http.Handler interface and routes requests to sub-handlers.
func (reg *Registry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := strings.TrimPrefix(r.URL.Path, "/admin")

	// 1. Static Asset Routing
	if strings.HasPrefix(upath, "/uploads/") {
		reg.handleStatic(w, r, upath)
		return
	}

	user, role := reg.GetUserFromRequest(r)

	// 2. Authentication Routing
	if upath == "/login" || upath == "/logout" {
		reg.routeAuth(w, r, upath)
		return
	}

	// 3. Auth Guard
	if user == nil {
		http.Redirect(w, r, "/admin/login", 303)
		return
	}

	// 4. Dashboard Routing
	if upath == "" || upath == "/" {
		reg.renderDashboard(w, r, user)
		return
	}

	// 5. Search API Routing
	if strings.HasSuffix(upath, "/search") {
		reg.routeSearch(w, r, upath)
		return
	}

	// 6. Main Resource/Page Routing
	reg.routeMain(w, r, upath, user, role)
}

func (reg *Registry) handleStatic(w http.ResponseWriter, r *http.Request, upath string) {
	fileName := strings.TrimPrefix(upath, "/uploads/")
	http.ServeFile(w, r, filepath.Join(reg.Config.UploadDir, fileName))
}

func (reg *Registry) routeAuth(w http.ResponseWriter, r *http.Request, upath string) {
	if upath == "/login" {
		if r.Method == "POST" {
			reg.handleLogin(w, r)
		} else {
			reg.renderLogin(w, r, "")
		}
		return
	}
	reg.handleLogout(w, r)
}

func (reg *Registry) routeSearch(w http.ResponseWriter, r *http.Request, upath string) {
	parts := strings.Split(strings.TrimPrefix(upath, "/"), "/")
	reg.handleSearchAPI(parts[0], w, r)
}

func (reg *Registry) routeMain(w http.ResponseWriter, r *http.Request, upath string, user *models.AdminUser, role string) {
	parts := strings.Split(strings.TrimPrefix(upath, "/"), "/")
	resourceName := parts[0]

	// Check Custom Pages
	if page, ok := reg.Pages[resourceName]; ok {
		page.Handler(w, r)
		return
	}

	// Check Resources
	res, ok := reg.GetResource(resourceName)
	if !ok {
		http.NotFound(w, r)
		return
	}

	action := "list"
	if len(parts) > 1 && parts[1] != "" {
		action = parts[1]
	}

	// Permission Check
	if !reg.IsAllowed(role, resourceName, action) && 
	   action != "export" && !strings.Contains(action, "action") {
		http.Error(w, "Forbidden", 403)
		return
	}

	reg.handleResourceAction(res, action, w, r, user)
}

func (reg *Registry) handleResourceAction(res *resource.Resource, action string, w http.ResponseWriter, r *http.Request, user *models.AdminUser) {
	switch action {
	case "export":
		reg.handleExport(res, w, r)
	case "action":
		reg.handleCustomAction(res, w, r, false)
	case "collection_action":
		reg.handleCustomAction(res, w, r, true)
	case "batch_action":
		reg.handleBatchAction(res, w, r)
	case "save":
		reg.handleSave(res, w, r, user)
	case "new":
		reg.renderForm(res, nil, w, r, user)
	case "show":
		id := r.URL.Query().Get("id")
		item, _ := reg.Get(res.Name, id)
		reg.renderShow(res, item, w, r, user)
	case "edit":
		id := r.URL.Query().Get("id")
		item, _ := reg.Get(res.Name, id)
		reg.renderForm(res, item, w, r, user)
	case "delete":
		id := r.URL.Query().Get("id")
		reg.Delete(res.Name, id)
		reg.RecordAction(user, res.Name, id, "Delete", "Record deleted")
		reg.setFlash(w, fmt.Sprintf("%s deleted successfully", res.Name))
		http.Redirect(w, r, "/admin/"+res.Name, 303)
	default:
		reg.renderList(res, w, r, user)
	}
}
