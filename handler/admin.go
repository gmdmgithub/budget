package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// PermissionType - permissions
type PermissionType struct {
	name string
}

// IsAdmin - first implementation
func (p PermissionType) IsAdmin() bool {
	return p.name == "admin"
}

// AdminRouter - a completely separate router for administrator routes
func AdminRouter() http.Handler {
	r := chi.NewRouter()
	// r.Use(AdminOnly)
	r.Get("/", adminIndex)
	r.Get("/accounts", adminListAccounts)
	return r
}

func adminIndex(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("Admin panel")))
}

func adminListAccounts(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("Admin list")))
}

// AdminOnly - check admin permission
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		perm, ok := ctx.Value("permission").(PermissionType)
		if !ok || !perm.IsAdmin() {
			http.Error(w, http.StatusText(403), 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}
