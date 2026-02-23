package handlers

import (
	"encoding/csv"
	"fmt"
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/resource"
	"net/http"
	"reflect"
)

func HandleExport(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s_export.csv", res.Name))
	writer := csv.NewWriter(w); defer writer.Flush()
	var h []string; for _, f := range res.Fields { h = append(h, f.Label) }; writer.Write(h)
	query := reg.DB.Model(res.Model)
	modelType := reflect.TypeOf(res.Model)
	destSlice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0); dest := reflect.New(destSlice.Type())
	query.Find(dest.Interface()); items := dest.Elem()
	for i := 0; i < items.Len(); i++ {
		item := reflect.Indirect(items.Index(i)); var row []string
		for _, f := range res.Fields { row = append(row, fmt.Sprintf("%v", item.FieldByName(f.Name).Interface())) }
		writer.Write(row)
	}
}
