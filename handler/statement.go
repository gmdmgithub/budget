package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// StatementRouter - a completely separate router for administrator routes
func StatementRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", allStatements)
	r.Route("/{stID}", func(r chi.Router) {
		r.Use(StatementCtx)
		r.Get("/", getStatement)       // GET /stat/123
		r.Put("/", updateStatement)    // PUT /stat/123
		r.Delete("/", deleteStatement) // DELETE /statement/123
		// Regexp url parameters:
		r.Get("/{name:[a-z-]+}", allStatements) // GET /statement/income from ABC
	})

	return r
}

// StatementCtx add Statement to the context
func StatementCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stID := chi.URLParam(r, "stID")
		//   statement, err := dbGetStatement(stID)
		//   if err != nil {
		// 	http.Error(w, http.StatusText(404), 404)
		// 	return
		//   }
		//   ctx := context.WithValue(r.Context(), "article", statement)
		ctx := context.WithValue(r.Context(), "stID", stID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func allStatements(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("All statements here!")))
}

func getStatement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// statement, ok := ctx.Value("statement").(*models.Statement)
	statement, ok := ctx.Value("stID").(string)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write([]byte(fmt.Sprintf("Specific statement %+v", statement)))

}
func updateStatement(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("Update statement?")))
}
func deleteStatement(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Delete statement?")))
}
