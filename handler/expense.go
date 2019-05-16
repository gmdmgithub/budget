package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gmdmgithub/budget/driver"
	"github.com/gmdmgithub/budget/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi"
)

// ExpensesRouter - router for expenses model
func ExpensesRouter() http.Handler {

	r := chi.NewRouter()

	r.Get("/", expenses)
	r.Post("/", createExpese)

	r.Route("/{ID}", func(r chi.Router) {
		r.Use(expenseCtx)
		r.Get("/", expense)
		r.Put("/", updateExpense)
		r.Delete("/", deleteExpense)
	})

	return r
}

func expenseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "ID")
		var exp model.Expense
		err := driver.DoOne(&exp, id, driver.GetOne)
		if err != nil {
			log.Printf("Problem - no expense in context %s %v", id, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "expense", &exp)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func expenses(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there all expenses here"))
}

func expense(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exp, ok := ctx.Value("expense").(*model.Expense)
	if !ok {
		log.Printf("Problem - no expense in context")
		http.Error(w, "Problem no expense in context", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&exp); err != nil {
		log.Printf("Problem - no expense in context %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// w.Write([]byte("Hi there one expense here"))
}

func createExpese(w http.ResponseWriter, r *http.Request) {

	//REMARK - date,to be properly decode, should be in format as follow
	// "date": "2019-04-01T00:00:00+01:00"
	var exp model.Expense
	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		log.Printf("Problem with decoding body %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exp.Created = time.Now()
	exp.UsrCreated = "234" //TODO! Change to proper

	res, err := driver.Create(&exp)
	if err != nil {
		log.Printf("Problem with creation expenses %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exp.ID = res.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&exp); err != nil {
		log.Printf("Problem with encoding expense %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func updateExpense(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exp, ok := ctx.Value("expense").(*model.Expense)
	if !ok {
		log.Printf("Problem - no expense in context")
		http.Error(w, "Problem no expense in context", http.StatusInternalServerError)
		return
	}

	var expenseBody model.Expense
	err := json.NewDecoder(r.Body).Decode(&expenseBody)
	if err != nil {
		log.Printf("Problem - no expense in body")
		http.Error(w, "Problem no expense in body", http.StatusInternalServerError)
		return
	}

	expenseBody.ID = exp.ID
	expenseBody.Created = exp.Created
	expenseBody.UsrCreated = exp.UsrCreated
	expenseBody.Updated = time.Now()
	expenseBody.UsrUpdated = "12345" //! TODO correct in the future

	res, err := driver.UpdateOne(&expenseBody, expenseBody.ID)
	if err != nil {
		log.Printf("Problem - with update expense")
		http.Error(w, "Problem - with update expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem update ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func deleteExpense(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exp, ok := ctx.Value("expense").(*model.Expense)
	if !ok {
		log.Printf("Problem - no expense in context")
		http.Error(w, "Problem no expense in context", http.StatusInternalServerError)
		return
	}
	res, err := driver.DeleteOne(exp, exp.ID)
	if err != nil {
		log.Printf("Problem delete Expense ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("delete expense"))
}
