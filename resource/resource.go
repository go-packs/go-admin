// Package resource defines metadata and helpers for admin-managed resources.
package resource

import (
	"html/template"
	"net/http"
	"reflect"

	"gorm.io/gorm"
)

// ActionHandler handles member or collection actions for a resource.
type ActionHandler func(res *Resource, w http.ResponseWriter, r *http.Request)

// BatchActionHandler handles batch actions over multiple selected IDs.
type BatchActionHandler func(res *Resource, ids []string, w http.ResponseWriter, r *http.Request)

// ScopeFunc filters a GORM query for a resource scope.
type ScopeFunc func(db *gorm.DB) *gorm.DB

// DecoratorFunc customizes field rendering.
type DecoratorFunc func(val interface{}) template.HTML

// SidebarHandler renders sidebar content for a resource item.
type SidebarHandler func(res *Resource, item interface{}) template.HTML

// Action represents a member or collection action.
type Action struct {
	Name, Label string
	Handler     ActionHandler
}

// BatchAction represents an action performed on multiple records.
type BatchAction struct {
	Name, Label string
	Handler     BatchActionHandler
}

// Scope represents a predefined filter/tab for a resource.
type Scope struct {
	Name, Label string
	Handler     ScopeFunc
}

// Sidebar represents a sidebar widget for a resource item.
type Sidebar struct {
	Label   string
	Handler SidebarHandler
}

// Association defines a relation between resources.
type Association struct{ Type, Name, ResourceName, ForeignKey, Label string }

// Field describes a single model field exposed in the admin UI.
type Field struct {
	Name, Label, Type string
	Options           []string
	Readonly          bool
	Searchable        bool
	SearchResource    string
	Decorator         DecoratorFunc
	Sortable          bool
}

// Resource holds metadata for a GORM model managed by the admin.
type Resource struct {
	Model             interface{}
	Name, Path, Group string
	Fields            []Field
	IndexFields       []string
	ShowFields        []string
	EditFields        []string
	MemberActions     []Action
	CollectionActions []Action
	BatchActions      []BatchAction
	Scopes            []Scope
	Associations      []Association
	Sidebars          []Sidebar
	Attributes        map[string]interface{}
}

// NewResource creates a new Resource metadata object from a model value.
func NewResource(model interface{}) *Resource {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return &Resource{Model: model, Name: t.Name(), Path: "/" + t.Name()}
}

func (r *Resource) SetGroup(group string) *Resource { r.Group = group; return r }
func (r *Resource) RegisterField(name, label string, readonly bool) *Resource {
	r.Fields = append(r.Fields, Field{Name: name, Label: label, Type: "text", Readonly: readonly, Sortable: true})
	return r
}
func (r *Resource) SetSortable(name string, sortable bool) *Resource {
	for i, f := range r.Fields {
		if f.Name == name {
			r.Fields[i].Sortable = sortable
			break
		}
	}
	return r
}
func (r *Resource) SetDecorator(name string, fn DecoratorFunc) *Resource {
	for i, f := range r.Fields {
		if f.Name == name {
			r.Fields[i].Decorator = fn
			break
		}
	}
	return r
}
func (r *Resource) AddSidebar(label string, handler SidebarHandler) *Resource {
	r.Sidebars = append(r.Sidebars, Sidebar{Label: label, Handler: handler})
	return r
}
func (r *Resource) AddMemberAction(n, l string, h ActionHandler) *Resource {
	r.MemberActions = append(r.MemberActions, Action{Name: n, Label: l, Handler: h})
	return r
}
func (r *Resource) AddCollectionAction(n, l string, h ActionHandler) *Resource {
	r.CollectionActions = append(r.CollectionActions, Action{Name: n, Label: l, Handler: h})
	return r
}
func (r *Resource) AddBatchAction(n, l string, h BatchActionHandler) *Resource {
	r.BatchActions = append(r.BatchActions, BatchAction{Name: n, Label: l, Handler: h})
	return r
}
func (r *Resource) AddScope(n, l string, h ScopeFunc) *Resource {
	r.Scopes = append(r.Scopes, Scope{Name: n, Label: l, Handler: h})
	return r
}
func (r *Resource) HasMany(n, l, tr, fk string) *Resource {
	r.Associations = append(r.Associations, Association{Type: "HasMany", Name: n, Label: l, ResourceName: tr, ForeignKey: fk})
	return r
}
func (r *Resource) BelongsTo(n, l, tr, fk string) *Resource {
	r.Associations = append(r.Associations, Association{Type: "BelongsTo", Name: n, Label: l, ResourceName: tr, ForeignKey: fk})
	return r
}
func (r *Resource) SetSearchable(f, tr string) *Resource {
	for i, field := range r.Fields {
		if field.Name == f {
			r.Fields[i].Searchable, r.Fields[i].SearchResource = true, tr
			break
		}
	}
	return r
}
func (r *Resource) SetFieldType(n, t string, opt ...string) *Resource {
	for i, f := range r.Fields {
		if f.Name == n {
			r.Fields[i].Type, r.Fields[i].Options = t, opt
			break
		}
	}
	return r
}
func (r *Resource) SetIndexFields(n ...string) *Resource { r.IndexFields = n; return r }
func (r *Resource) SetShowFields(n ...string) *Resource  { r.ShowFields = n; return r }
func (r *Resource) SetEditFields(n ...string) *Resource  { r.EditFields = n; return r }

func (r *Resource) GetFieldsFor(view string) []Field {
	var names []string
	switch view {
	case "index":
		names = r.IndexFields
	case "show":
		names = r.ShowFields
	case "edit", "new":
		names = r.EditFields
	}
	if len(names) == 0 {
		return r.Fields
	}
	var result []Field
	for _, name := range names {
		for _, f := range r.Fields {
			if f.Name == name {
				result = append(result, f)
				break
			}
		}
	}
	return result
}
