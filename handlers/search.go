package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-packs/go-admin"
	"net/http"
	"reflect"
	"strings"
)

func HandleSearchAPI(reg *admin.Registry, resourceName string, w http.ResponseWriter, r *http.Request) {
	res, ok := reg.GetResource(resourceName); if !ok { http.Error(w, "Not found", 404); return }
	query := r.URL.Query().Get("q"); db := reg.DB.Model(res.Model); searchQuery := ""
	for _, f := range res.Fields { if f.Type == "text" { if searchQuery != "" { searchQuery += " OR " }; searchQuery += fmt.Sprintf("%s LIKE ?", f.Name) } }
	if searchQuery != "" { args := make([]interface{}, strings.Count(searchQuery, "?")); for i := range args { args[i] = "%" + query + "%" }; db = db.Where(searchQuery, args...) }
	var results []map[string]interface{}; modelType := reflect.TypeOf(res.Model)
	destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0); dest := reflect.New(destSlice.Type())
	db.Limit(10).Find(dest.Interface()); items := dest.Elem()
	for i := 0; i < items.Len(); i++ {
		item := reflect.Indirect(items.Index(i)); m := make(map[string]interface{})
		m["id"] = item.FieldByName("ID").Interface()
		if f := item.FieldByName("Name"); f.IsValid() { m["text"] = f.Interface() } else if f := item.FieldByName("Email"); f.IsValid() { m["text"] = f.Interface() } else { m["text"] = fmt.Sprintf("ID: %v", m["id"]) }
		results = append(results, m)
	}
	w.Header().Set("Content-Type", "application/json"); json.NewEncoder(w).Encode(results)
}
