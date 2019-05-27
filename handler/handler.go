package handler

import "net/http"

type contextKey string

// GlobalCtx - func for global middleware - like content-type
func GlobalCtx(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		//add more if necessary
		next.ServeHTTP(w, r)
	})
}
