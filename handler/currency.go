package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gmdmgithub/budget/driver"
	"github.com/gmdmgithub/budget/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi"
)

var cL []string

type checkTime struct {
	time  time.Time
	valid bool // Valid is true if Time is not NULL
}

var cT = checkTime{
	valid: false,
}

// CurrencyRouter - group all routes for currency model
func CurrencyRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(currencyCxt)
	r.Get("/", currencies)
	r.Post("/", createCurrency) //create one currency
	r.Route("/{curCode}", func(r chi.Router) {
		r.Get("/", currency)
		r.Get("/date", currencyDate) //GET currency from specific date (nearest), default current
	})

	return r
}

func currencyCxt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//at least once a day check what currencies are stored
		if !cT.valid {
			cT.time = time.Now()
		}
		if time.Since(cT.time) > (24*time.Hour) || !cT.valid {
			cT.time = time.Now()
			cT.valid = true
			res, err := driver.Distinct("currencies", "code")
			if err != nil {
				log.Printf("Problem with currency codes %v", err)
			}
			if len(res) > 0 {
				cL = nil
			}
			for _, c := range res {
				cL = append(cL, fmt.Sprintf("%s", c))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func createCurrency(w http.ResponseWriter, r *http.Request) {

	var cur model.Currency

	if err := json.NewDecoder(r.Body).Decode(&cur); err != nil {
		log.Printf("Problem with decode currency %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var m model.Modeler = &cur

	log.Printf("Let's see the model %+v", m)

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

func currencyDate(w http.ResponseWriter, r *http.Request) {

	log.Printf("URL queries %+v", r.URL.Query())

	curCode := chi.URLParam(r, "curCode")
	if curCode != "USD" && curCode != "EUR" && curCode != "CHF" && curCode != "GBP" {
		curCode = "PLN"
	}

	// options
	opts := options.Find()
	opts.SetLimit(1)
	// Sort by `date` field +1 ascending -1 descending
	opts.Sort = bson.M{"date": 1}

	params := r.URL.Query()
	layout := "2006-01-02"
	s := params.Get("date")
	t, err := time.Parse(layout, s)
	if err != nil {
		log.Printf("Problem with parsing date so current date is taken: %v", err)
		t = time.Now()
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}
	tto := time.Now().Add(24 * time.Hour)

	filter := bson.M{}
	filter["date"] = bson.M{"$gte": t, "$lte": tto}
	filter["code"] = bson.M{"$eq": curCode}

	// proper db query
	curs, err := driver.GetCurrencies(filter, opts)
	if err != nil {
		log.Printf("Problem with get currencies %v", curs)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// simplify response - object instead array
	var cur model.Currency
	if len(curs) == 1 {
		cur = curs[0]
	} else {
		cur = model.Currency{
			Code:         "PLN",
			ExchangeRate: 1,
			Date:         time.Now(),
			Base:         true,
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&cur); err != nil {
		log.Printf("Problem with encoding currencies %v", cur)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// db.collection.find({"time":{$gte: isoDate,$lt: isoDate}}).sort({"time":1}).limit(1)

	// fmt.Fprintf(w, "Hi there - purposefully instead of w.Write([]byte(\"Hi there currency for one date here\"))")
}

// currencies - get list of all or lats
func currencies(w http.ResponseWriter, r *http.Request) {

	opt := options.Find()

	filter := bson.M{}

	recent := r.URL.Query().Get("recent")
	log.Printf("give me recent %v", recent)
	if recent == "true" {
		// TODO build filter for the last exchange rate for each currency
		log.Printf("Recent is true %v", recent)
		// filter
	}

	cur, err := driver.GetCurrencies(filter, opt)
	if err != nil {
		log.Printf("Problem with get currencies %v", cur)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(cur); err != nil {
		log.Printf("Problem with encoding currencies %v", cur)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("Hi there all currencies here"))
}

func currency(w http.ResponseWriter, r *http.Request) {

	// similar to date - the nearest previous currency - to be rewrite - first simple implementation - unfortunately not follow DRY
	curCode := chi.URLParam(r, "curCode")
	if curCode != "USD" && curCode != "EUR" && curCode != "CHF" && curCode != "GBP" {
		curCode = "PLN"
	}

	// options
	opts := options.Find()
	opts.SetLimit(1)
	// Sort by `date` field +1 ascending -1 descending
	opts.Sort = bson.M{"date": -1}

	tto := time.Now().Add(24 * time.Hour)

	filter := bson.M{}
	filter["date"] = bson.M{"$lte": tto}
	filter["code"] = bson.M{"$eq": curCode}

	// proper db query
	curs, err := driver.GetCurrencies(filter, opts)
	if err != nil {
		log.Printf("Problem with get currencies %v", curs)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// simplify response - object instead array
	var cur model.Currency
	if len(curs) == 1 {
		cur = curs[0]
	} else {
		cur = model.Currency{
			Code:         "PLN",
			ExchangeRate: 1,
			Date:         time.Now(),
			Base:         true,
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&cur); err != nil {
		log.Printf("Problem with encoding currencies %v", cur)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("Hi there one currency here"))
}
