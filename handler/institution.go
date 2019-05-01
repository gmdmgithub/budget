package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func InstitutionRouter() http.Handler {

	r := chi.NewRouter()

	r.Get("/", institutions()) //all- using
	r.Route("/{instID}", func(r chi.Router) {
		r.Get("/", getInstitution) // GET /institution/123
	})

	return r
}

func institutions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(fmt.Sprintf("All institutions here!")))
	}
}

func getInstitution(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("One institution here!")))
}
