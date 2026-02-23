// Package admin provides a lightweight admin panel registry and helpers.
package admin

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-packs/go-admin/config"
	"github.com/go-packs/go-admin/models"
	"github.com/go-packs/go-admin/resource"
	"gorm.io/gorm"
)

//go:embed templates/*
var TemplateFS embed.FS

// Public Type Aliases
// Resource is an alias for resource.Resource.
type Resource = resource.Resource

// Field is an alias for resource.Field.
type Field = resource.Field

// Config is an alias for config.Config.
type Config = config.Config

// AdminUser is an alias for models.AdminUser.
type AdminUser = models.AdminUser

// Session is an alias for models.Session.
type Session = models.Session

// Permission is an alias for models.Permission.
type Permission = models.Permission

// AuditLog is an alias for models.AuditLog.
type AuditLog = models.AuditLog

// Scope is an alias for resource.Scope.
type Scope = resource.Scope

// Registry coordinates registered resources, pages, charts and configuration.
type Registry struct {
	DB        *gorm.DB
	Resources map[string]*resource.Resource
	Pages     map[string]*Page
	Charts    []Chart
	Config    *config.Config
}

// Page describes a custom admin page.
type Page struct {
	Name, Group string
	Handler     http.HandlerFunc
}

// Chart represents a dashboard chart widget.
type Chart struct {
	Label string
	Type  string
	Data  func(db *gorm.DB) (labels []string, values []float64)
}

// NewRegistry creates a new admin Registry.
func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		DB: db, Resources: make(map[string]*resource.Resource), Pages: make(map[string]*Page),
		Charts: []Chart{}, Config: config.DefaultConfig(),
	}
}

// Public Factory Functions
// DefaultConfig returns the default configuration.
func DefaultConfig() *config.Config { return config.DefaultConfig() }

// LoadConfig loads configuration from a file.
func LoadConfig(path string) (*config.Config, error) { return config.LoadConfig(path) }

// NewResource creates a resource from a model.
func NewResource(model interface{}) *resource.Resource { return resource.NewResource(model) }

func (reg *Registry) SetConfig(c *config.Config) { reg.Config = c }

func (reg *Registry) AddChart(l, t string, p func(db *gorm.DB) ([]string, []float64)) {
	reg.Charts = append(reg.Charts, Chart{Label: l, Type: t, Data: p})
}

func (reg *Registry) AddPage(n, g string, h http.HandlerFunc) {
	reg.Pages[n] = &Page{Name: n, Group: g, Handler: h}
}

func (reg *Registry) Register(m interface{}) *resource.Resource {
	res := resource.NewResource(m)
	reg.Resources[res.Name] = res
	fmt.Printf("Registered resource: %s\n", res.Name)
	return res
}

func (reg *Registry) GetResource(n string) (*resource.Resource, bool) {
	res, ok := reg.Resources[n]
	return res, ok
}

func (reg *Registry) ResourceNames() []string {
	names := make([]string, 0, len(reg.Resources))
	for n := range reg.Resources {
		names = append(names, n)
	}
	return names
}

func (reg *Registry) GetGroupedResources() map[string][]*resource.Resource {
	groups := make(map[string][]*resource.Resource)
	for _, r := range reg.Resources {
		g := r.Group
		if g == "" {
			g = "Default"
		}
		groups[g] = append(groups[g], r)
	}
	return groups
}

func (reg *Registry) GetGroupedPages() map[string][]*Page {
	groups := make(map[string][]*Page)
	for _, p := range reg.Pages {
		g := p.Group
		if g == "" {
			g = "Default"
		}
		groups[g] = append(groups[g], p)
	}
	return groups
}

func (reg *Registry) SetFlash(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{Name: "admin_flash", Value: message, Path: "/admin", HttpOnly: true})
}

func (reg *Registry) GetFlash(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("admin_flash")
	if err != nil {
		return ""
	}
	http.SetCookie(w, &http.Cookie{Name: "admin_flash", Value: "", Path: "/admin", MaxAge: -1})
	return cookie.Value
}
