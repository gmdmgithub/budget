package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gmdmgithub/budget/driver"
	"github.com/gmdmgithub/budget/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi"
)

func CurrencyRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", currencies)
	r.Post("/", createCurrency) //create one currency
	r.Route("/{curCode}", func(r chi.Router) {
		r.Get("/", currency)
	})

	return r
}

func createCurrency(w http.ResponseWriter, r *http.Request) {

	var cur model.Currency

	if err := json.NewDecoder(r.Body).Decode(&cur); err != nil {
		log.Printf("Problem with decode curr %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var m model.Modeler = &cur

	res, err := driver.Create(m)
	if err != nil {
		log.Printf("Saving problem %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cur.ID = res.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(cur); err != nil {
		log.Printf("Encode problem %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func currencies(w http.ResponseWriter, r *http.Request) {

	cur, err := driver.GetAllCurrencies()
	if err != nil {
		log.Printf("Problem with get currencies %v", cur)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cur); err != nil {
		log.Printf("Problem with encoding currencies %v", cur)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("Hi there all currencies here"))
}

func currency(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there one currency here"))
}
