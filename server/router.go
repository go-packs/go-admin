package server

import (
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/handlers"
	"github.com/go-packs/go-admin/internal"
	"github.com/go-packs/go-admin/view"
	"net/http"
	"path/filepath"
	"strings"
)

func NewRouter(reg *admin.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := strings.TrimPrefix(r.URL.Path, "/admin")

		// 1. Static Asset Routing
		if strings.HasPrefix(upath, "/uploads/") {
			handleStatic(reg, w, r, upath)
			return
		}

		user, role := internal.GetUserFromRequest(reg, r)

		// 2. Authentication Routing
		if upath == "/login" || upath == "/logout" {
			routeAuth(reg, w, r, upath)
			return
		}

		// 3. Auth Guard
		if user == nil {
			http.Redirect(w, r, "/admin/login", 303)
			return
		}

		// 4. Dashboard Routing
		if upath == "" || upath == "/" {
			view.RenderDashboard(reg, w, r, user)
			return
		}

		// 5. Search API Routing
		if strings.HasSuffix(upath, "/search") {
			parts := strings.Split(strings.TrimPrefix(upath, "/"), "/")
			handlers.HandleSearchAPI(reg, parts[0], w, r)
			return
		}

		// 6. Main Resource/Page Routing
		routeMain(reg, w, r, upath, user, role)
	})
}

func handleStatic(reg *admin.Registry, w http.ResponseWriter, r *http.Request, upath string) {
	fileName := strings.TrimPrefix(upath, "/uploads/")
	http.ServeFile(w, r, filepath.Join(reg.Config.UploadDir, fileName))
}

func routeAuth(reg *admin.Registry, w http.ResponseWriter, r *http.Request, upath string) {
	if upath == "/login" {
		if r.Method == "POST" {
			handlers.Login(reg)(w, r)
		} else {
			handlers.RenderLogin(reg, w, r, "")
		}
		return
	}
	handlers.Logout(reg)(w, r)
}

func routeMain(reg *admin.Registry, w http.ResponseWriter, r *http.Request, upath string, user *admin.AdminUser, role string) {
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
	if !internal.IsAllowed(reg, role, resourceName, action) && 
	   action != "export" && !strings.Contains(action, "action") {
		http.Error(w, "Forbidden", 403)
		return
	}

	handleResourceAction(reg, res, action, w, r, user)
}

func handleResourceAction(reg *admin.Registry, res *admin.Resource, action string, w http.ResponseWriter, r *http.Request, user *admin.AdminUser) {
	switch action {
	case "export":
		handlers.HandleExport(reg, res, w, r)
	case "action":
		handlers.HandleCustomAction(reg, res, w, r, false)
	case "collection_action":
		handlers.HandleCustomAction(reg, res, w, r, true)
	case "batch_action":
		handlers.HandleBatchAction(reg, res, w, r)
	case "save":
		handlers.HandleSave(reg, res, w, r, user)
	case "new":
		handlers.RenderForm(reg, res, nil, w, r, user)
	case "show":
		id := r.URL.Query().Get("id")
		item, _ := internal.Get(reg, res.Name, id)
		handlers.RenderShow(reg, res, item, w, r, user)
	case "edit":
		id := r.URL.Query().Get("id")
		item, _ := internal.Get(reg, res.Name, id)
		handlers.RenderForm(reg, res, item, w, r, user)
	case "delete":
		handlers.HandleDelete(reg, res, w, r, user)
	default:
		handlers.RenderList(reg, res, w, r, user)
	}
}
