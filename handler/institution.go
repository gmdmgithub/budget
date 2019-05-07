package handler

import (
	"encoding/json"
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

// InstitutionRouter main router for institutions
func InstitutionRouter() http.Handler {

	r := chi.NewRouter()

	r.Get("/", institutions())     //all- using a function - like middleware way
	r.Post("/", createInstitution) //
	r.Route("/{code}", func(r chi.Router) {
		r.Get("/", getInstitution) // GET /institution/123
	})

	return r
}

func createInstitution(w http.ResponseWriter, r *http.Request) {

	var i model.Institution

	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		log.Printf("Problem with decoding Institution: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i.Created = time.Now()
	i.UsrCreated = "23" //change
	var m model.Modeler = &i

	res, err := driver.Create(m)
	if err != nil {
		log.Printf("Problem with creating Institution: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i.ID = res.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(i); err != nil {
		log.Printf("Problem with encoding Institution: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("Create institution is here!"))
}

func institutions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		filer := bson.M{}
		opt := options.Find()

		var i model.Institution

		res, err := driver.GetList(filer, opt, &i)
		if err != nil {
			log.Printf("Problem with List of Institution: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var insts []model.Institution
		for _, inst := range res {
			inst, ok := inst.(*model.Institution)
			if ok {
				insts = append(insts, *inst)
			}
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&insts); err != nil {
			log.Printf("Problem with encoding Institutions %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// w.Write([]byte(fmt.Sprintf("All institutions here!")))
	}
}

func getInstitution(w http.ResponseWriter, r *http.Request) {

	filter := bson.M{}
	opt := options.FindOne()

	// filter["code"] = bson.M{"$regex": `(?i)` + chi.URLParam(r, "code")}
	// log.Printf("Filter is %+v", filter["code"])
	filter["code"] = chi.URLParam(r, "code")

	var i model.Institution

	err := driver.GetOne(&i, filter, opt)
	if err != nil {
		log.Printf("Problem with Institution %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&i); err != nil {
		log.Printf("Problem with encoding Institution %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte(fmt.Sprintf("One institution here!")))
}
