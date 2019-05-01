package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func StmntTypeRouter() http.HandlerFunc {
	r := chi.NewRouter()

	return r
}
