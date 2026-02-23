package handlers

import (
	"net/http"

	"github.com/go-packs/go-admin"
	"github.com/go-packs/go-admin/internal"
)

func AuthGuard(reg *admin.Registry, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := internal.GetUserFromRequest(reg, r)
		if user == nil {
			http.Redirect(w, r, "/admin/login", 303)
			return
		}
		next(w, r)
	}
}
