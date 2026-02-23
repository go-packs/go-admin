package handlers

import (
	"fmt"
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/internal"
	"github.com/go-packs/go-admin/models"
	"github.com/go-packs/go-admin/resource"
	"github.com/go-packs/go-admin/view"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func RenderList(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request, user *models.AdminUser) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fields := res.GetFieldsFor("index")
	page, _ := strconv.Atoi(r.URL.Query().Get("page")); if page < 1 { page = 1 }
	perPage := reg.Config.DefaultPerPage
	currentScope := r.URL.Query().Get("scope")
	query := reg.DB.Model(res.Model)
	if currentScope != "" {
		for _, s := range res.Scopes { if s.Name == currentScope { query = s.Handler(query); break } }
	}
	sortField, sortOrder := r.URL.Query().Get("sort"), r.URL.Query().Get("order")
	if sortField != "" { if sortOrder != "desc" { sortOrder = "asc" }; query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder)) } else { query = query.Order("id desc") }
	filters := make(map[string]string)
	for k, v := range r.URL.Query() {
		val := v[0]; if val == "" { continue }; filters[k] = val
		if strings.HasPrefix(k, "q_") { query = query.Where(fmt.Sprintf("%s LIKE ?", strings.TrimPrefix(k, "q_")), "%"+val+"%") } else if strings.HasPrefix(k, "min_") { query = query.Where(fmt.Sprintf("%s >= ?", strings.TrimPrefix(k, "min_")), val) } else if strings.HasPrefix(k, "max_") { query = query.Where(fmt.Sprintf("%s <= ?", strings.TrimPrefix(k, "max_")), val) }
	}
	var totalCount int64; query.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	modelType := reflect.TypeOf(res.Model)
	destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0); dest := reflect.New(destSlice.Type())
	query.Offset((page - 1) * perPage).Limit(perPage).Find(dest.Interface())
	data := view.SliceToMap(res, fields, dest.Elem())
	styleContent, _ := admin.TemplateFS.ReadFile("templates/style.css")
	tmpl := view.LoadTemplates("templates/index.html")
	pd := view.PageData{
		SiteTitle: reg.Config.SiteTitle, Resources: reg.Resources, GroupedResources: reg.GetGroupedResources(), GroupedPages: reg.GetGroupedPages(), 
		CurrentResource: res, Fields: fields, Data: data, Filters: filters, User: user, CSS: template.CSS(styleContent),
		Page: page, PerPage: perPage, TotalPages: totalPages, TotalCount: totalCount, HasPrev: page > 1, HasNext: page < totalPages, PrevPage: page - 1, NextPage: page + 1, Scopes: res.Scopes, CurrentScope: currentScope,
		Flash: reg.GetFlash(w, r), SortField: sortField, SortOrder: sortOrder,
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", pd); err != nil {
		fmt.Printf("Template error: %v\n", err)
	}
}

func RenderShow(reg *admin.Registry, res *resource.Resource, item interface{}, w http.ResponseWriter, r *http.Request, user *models.AdminUser) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fields := res.GetFieldsFor("show")
	var itemMap map[string]interface{}
	assocData := make(map[string]*view.AssociationData); renderedSidebars := make(map[string]template.HTML)
	if item != nil {
		itemMap = view.ItemToMap(res, fields, reflect.ValueOf(item))
		for _, assoc := range res.Associations {
			if assoc.Type == "HasMany" {
				targetRes, _ := reg.GetResource(assoc.ResourceName)
				targetFields := targetRes.GetFieldsFor("index")
				modelType := reflect.TypeOf(targetRes.Model)
				destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0); dest := reflect.New(destSlice.Type())
				reg.DB.Where(fmt.Sprintf("%s = ?", assoc.ForeignKey), itemMap["ID"]).Find(dest.Interface())
				assocData[assoc.Name] = &view.AssociationData{Resource: targetRes, Fields: targetFields, Items: view.SliceToMap(targetRes, targetFields, dest.Elem())}
			}
		}
		for _, sb := range res.Sidebars { renderedSidebars[sb.Label] = sb.Handler(res, item) }
	}
	styleContent, _ := admin.TemplateFS.ReadFile("templates/style.css")
	tmpl := view.LoadTemplates("templates/show.html")
	pd := view.PageData{SiteTitle: reg.Config.SiteTitle, Resources: reg.Resources, GroupedResources: reg.GetGroupedResources(), GroupedPages: reg.GetGroupedPages(), CurrentResource: res, Fields: fields, Item: itemMap, User: user, CSS: template.CSS(styleContent), Associations: assocData, Flash: reg.GetFlash(w, r), RenderedSidebars: renderedSidebars}
	if err := tmpl.ExecuteTemplate(w, "show.html", pd); err != nil {
		fmt.Printf("Template error: %v\n", err)
	}
}

func RenderForm(reg *admin.Registry, res *resource.Resource, item interface{}, w http.ResponseWriter, r *http.Request, user *models.AdminUser) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	viewType := "edit"
	if item == nil {
		viewType = "new"
	}
	fields := res.GetFieldsFor(viewType)
	
	var itemMap map[string]interface{}
	if item != nil { itemMap = view.ItemToMap(res, fields, reflect.ValueOf(item)) }
	assocData := make(map[string]*view.AssociationData)
	for _, assoc := range res.Associations {
		if assoc.Type == "BelongsTo" {
			targetRes, _ := reg.GetResource(assoc.ResourceName)
			var count int64; reg.DB.Model(targetRes.Model).Count(&count)
			if count < reg.Config.SearchThreshold {
				modelType := reflect.TypeOf(targetRes.Model)
				destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0); dest := reflect.New(destSlice.Type())
				reg.DB.Find(dest.Interface())
				assocData[assoc.Name] = &view.AssociationData{Resource: targetRes, Options: view.SliceToMap(targetRes, targetRes.Fields, dest.Elem())}
			} else { assocData[assoc.Name] = &view.AssociationData{Resource: targetRes} }
		}
	}
	for _, f := range fields { if f.Searchable && f.SearchResource != "" { targetRes, _ := reg.GetResource(f.SearchResource); assocData[f.Name] = &view.AssociationData{Resource: targetRes} } }
	styleContent, _ := admin.TemplateFS.ReadFile("templates/style.css")
	tmpl := view.LoadTemplates("templates/form.html")
	pd := view.PageData{SiteTitle: reg.Config.SiteTitle, Resources: reg.Resources, GroupedResources: reg.GetGroupedResources(), GroupedPages: reg.GetGroupedPages(), CurrentResource: res, Fields: fields, Item: itemMap, User: user, CSS: template.CSS(styleContent), Associations: assocData, Flash: reg.GetFlash(w, r)}
	if err := tmpl.ExecuteTemplate(w, "form.html", pd); err != nil {
		fmt.Printf("Template error: %v\n", err)
	}
}

func HandleSave(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request, user *models.AdminUser) {
	r.ParseMultipartForm(32 << 20)
	model := reflect.New(reflect.TypeOf(res.Model)).Interface()
	isUpdate, id := false, r.FormValue("ID")
	if id != "" && id != "0" {
		reg.DB.First(model, id)
		isUpdate = true
	}

	elem := reflect.ValueOf(model).Elem()

	// Always try to set ID if it exists and we have a value, even if it's not in res.Fields
	if id != "" && id != "0" {
		idField := elem.FieldByName("ID")
		if idField.IsValid() && idField.CanSet() {
			if idField.Kind() == reflect.Uint {
				uv, _ := strconv.ParseUint(id, 10, 64)
				idField.SetUint(uv)
			} else if idField.Kind() == reflect.Int || idField.Kind() == reflect.Int64 {
				iv, _ := strconv.ParseInt(id, 10, 64)
				idField.SetInt(iv)
			}
		}
	}

	for _, f := range res.Fields {
		if f.Readonly {
			continue
		}
		field := elem.FieldByName(f.Name)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		if f.Type == "image" || f.Type == "file" {
			file, header, err := r.FormFile(f.Name)
			if err == nil {
				defer file.Close()
				os.MkdirAll(reg.Config.UploadDir, 0755)
				newName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
				dst, _ := os.Create(filepath.Join(reg.Config.UploadDir, newName))
				defer dst.Close()
				io.Copy(dst, file)
				field.SetString("/admin/uploads/" + newName)
			}
			continue
		}
		val := r.FormValue(f.Name)
		switch field.Kind() {
		case reflect.String:
			field.SetString(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uv, _ := strconv.ParseUint(val, 10, 64)
			field.SetUint(uv)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			iv, _ := strconv.ParseInt(val, 10, 64)
			field.SetInt(iv)
		case reflect.Float32, reflect.Float64:
			fv, _ := strconv.ParseFloat(val, 64)
			field.SetFloat(fv)
		case reflect.Bool:
			field.SetBool(val == "true" || val == "on" || val == "1")
		}
	}
	reg.DB.Save(model)
	newID := fmt.Sprintf("%v", elem.FieldByName("ID").Interface())
	act := "Create"
	if isUpdate {
		act = "Update"
	}
	internal.RecordAction(reg, user, res.Name, newID, act, "Saved from form")
	reg.SetFlash(w, fmt.Sprintf("%s saved successfully", res.Name))
	http.Redirect(w, r, "/admin/"+res.Name, 303)
}

func HandleDelete(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request, user *models.AdminUser) {
	id := r.URL.Query().Get("id")
	internal.Delete(reg, res.Name, id)
	internal.RecordAction(reg, user, res.Name, id, "Delete", "Record deleted")
	reg.SetFlash(w, fmt.Sprintf("%s deleted successfully", res.Name))
	http.Redirect(w, r, "/admin/"+res.Name, 303)
}
