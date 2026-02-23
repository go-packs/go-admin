package handlers

import (
	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/resource"
	"net/http"
)

func HandleBatchAction(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { http.Error(w, "Method not allowed", 405); return }
	r.ParseForm(); actionName, ids := r.FormValue("action_name"), r.Form["ids"]
	if actionName == "" || len(ids) == 0 { http.Redirect(w, r, "/admin/"+res.Name, 303); return }
	for _, a := range res.BatchActions { if a.Name == actionName { a.Handler(res, ids, w, r); return } }
}

func HandleCustomAction(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request, isCollection bool) {
	actionName := r.URL.Query().Get("name")
	var actions []resource.Action
	if isCollection { actions = res.CollectionActions } else { actions = res.MemberActions }
	for _, a := range actions { if a.Name == actionName { a.Handler(res, w, r); return } }
}
