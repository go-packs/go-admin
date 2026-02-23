// Package handlers provides HTTP handlers for the admin panel.
package handlers

import (
	"net/http"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/resource"
)

func HandleBatchAction(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	actionName, ids := r.FormValue("action_name"), r.Form["ids"]
	if actionName == "" || len(ids) == 0 {
		http.Redirect(w, r, "/admin/"+res.Name, 303)
		return
	}
	for _, a := range res.BatchActions {
		if a.Name == actionName {
			a.Handler(res, ids, w, r)
			return
		}
	}
}

func HandleCustomAction(reg *admin.Registry, res *resource.Resource, w http.ResponseWriter, r *http.Request, isCollection bool) {
	actionName := r.URL.Query().Get("name")
	var actions []resource.Action
	if isCollection {
		actions = res.CollectionActions
	} else {
		actions = res.MemberActions
	}
	for _, a := range actions {
		if a.Name == actionName {
			a.Handler(res, w, r)
			return
		}
	}
}
