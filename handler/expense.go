package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

// ExpensesRouter - router for expenses model
func ExpensesRouter() http.Handler {

	r := chi.NewRouter()

	r.Get("/", expenses)
	r.Post("/", createExpese)

	r.Route("/{expID}", func(r chi.Router) {
		r.Use(expenseCtx)
		r.Get("/", expense)
		r.Put("/", updateExpense)
		r.Delete("/", deleteExpense)
	})

	return r
}

func expenseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//TODO! get expense here and put it into the context

		next.ServeHTTP(w, r)
	})
}

func expenses(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there all expenses here"))
}

func expense(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there one expense here"))
}

func createExpese(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create expense here!"))
}

func updateExpense(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PUT for expense here"))
}

func deleteExpense(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete expense"))
}
