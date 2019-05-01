package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func StmntTypeRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", stmntTypes)

	r.Route("/{code}", func(r chi.Router) {
		r.Get("/", stmntType)
	})
	return r
}

func stmntTypes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there all statement types here"))
}

func stmntType(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there one statement type here"))
}
