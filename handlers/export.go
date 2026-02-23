package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/resource"
)

func HandleExport(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s_export.csv", res.Name))
	writer := csv.NewWriter(w)
	defer writer.Flush()
	var h []string
	for _, f := range res.Fields {
		h = append(h, f.Label)
	}
	if err := writer.Write(h); err != nil {
		http.Error(w, "Internal error", 500)
		return
	}
	query := reg.DB.Model(res.Model)
	modelType := reflect.TypeOf(res.Model)
	destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)
	dest := reflect.New(destSlice.Type())
	if err := query.Find(dest.Interface()).Error; err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	items := dest.Elem()
	for i := 0; i < items.Len(); i++ {
		item := reflect.Indirect(items.Index(i))
		var row []string
		for _, f := range res.Fields {
			row = append(row, fmt.Sprintf("%v", item.FieldByName(f.Name).Interface()))
		}
		if err := writer.Write(row); err != nil {
			http.Error(w, "Internal error", 500)
			return
		}
	}
}
