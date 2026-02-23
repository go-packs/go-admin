package view

import (
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/models"
	"github.com/go-packs/go-admin/resource"
	"html/template"
	"reflect"
)

type PageData struct {
	SiteTitle        string
	Resources        map[string]*resource.Resource
	GroupedResources map[string][]*resource.Resource
	GroupedPages     map[string][]*admin.Page
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
	Associations     map[string]*AssociationData
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

func LoadTemplates(contentTmpl string) *template.Template {
	return template.Must(template.ParseFS(admin.TemplateFS, "templates/layout.html", contentTmpl))
}

func SliceToMap(res *resource.Resource, fields []resource.Field, slice reflect.Value) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < slice.Len(); i++ { data = append(data, ItemToMap(res, fields, slice.Index(i))) }
	return data
}

func ItemToMap(res *resource.Resource, fields []resource.Field, item reflect.Value) map[string]interface{} {
	m := make(map[string]interface{})
	item = reflect.Indirect(item)
	for _, f := range fields {
		fv := item.FieldByName(f.Name)
		if fv.IsValid() {
			val := fv.Interface()
			if f.Decorator != nil {
				m[f.Name] = f.Decorator(val)
			} else {
				m[f.Name] = val
			}
		}
	}
	idv := item.FieldByName("ID"); if idv.IsValid() { m["ID"] = idv.Interface() }
	return m
}
