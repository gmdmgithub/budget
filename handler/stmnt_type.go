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

func StmntTypeRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", stmntTypes)
	r.Post("/", createStmntType)
	r.Route("/{code}", func(r chi.Router) {
		r.Get("/", stmntType)
	})
	return r
}

func createStmntType(w http.ResponseWriter, r *http.Request) {

	var sType model.StmntType

	if err := json.NewDecoder(r.Body).Decode(&sType); err != nil {
		log.Printf("Problem with decoding stmntType %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sType.Created = time.Now()
	sType.UsrCreated = "324" //TODO change to propper

	res, err := driver.Create(&sType)
	if err != nil {
		log.Printf("Problem to create stmtType %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sType.ID = res.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sType); err != nil {
		log.Printf("Problem with encoding stmtType %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("Lets create a statement type"))
}

func stmntTypes(w http.ResponseWriter, r *http.Request) {

	filter := bson.M{}
	opt := options.Find()

	var stmntType model.StmntType
	res, err := driver.GetList(filter, opt, &stmntType)
	if err != nil {
		log.Printf("Problem with stmt type list %v", err)
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	var stts []model.StmntType
	for _, stt := range res {
		s, _ := stt.(*model.StmntType)
		stts = append(stts, *s)
	}
	if err := json.NewEncoder(w).Encode(stts); err != nil {
		log.Printf("Problem with encoding stmt type list %v", err)
		http.Error(w, err.Error(), http.StatusOK)
	}

	// w.Write([]byte("Hi there all statement types here"))
}

func stmntType(w http.ResponseWriter, r *http.Request) {

	filter := bson.M{}
	code := chi.URLParam(r, "code")
	filter["code"] = code
	var st model.StmntType
	opt := options.FindOne()

	err := driver.GetOne(&st, filter, opt)
	if err != nil {
		log.Printf("Get statement type problem %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&st); err != nil {
		log.Printf("Get statement type problem %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("Hi there one statement type here"))
}
