package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func CurrencyRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", currencies)

	r.Route("/{curCode}", func(r chi.Router) {
		r.Get("/", currency)
	})

	return r
}

func currencies(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there all currencies here"))
}

func currency(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there one currency here"))
}
