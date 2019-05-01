package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func ExpensesRouter() http.Handler {

	r := chi.NewRouter()

	r.Get("/", expenses)

	r.Route("/{expID}", func(r chi.Router) {
		r.Get("/", expense)

	})

	return r
}

func expenses(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there all expenses here"))
}

func expense(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there one expense here"))
}
