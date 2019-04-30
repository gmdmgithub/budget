package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gmdmgithub/budget/driver"
	"github.com/gmdmgithub/budget/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func UserRouter() http.Handler {

	r := chi.NewRouter()

	r.Get("/", getAllUsers)
	r.Post("/", createUser)
	r.Route("/{usrID}", func(r chi.Router) {
		r.Get("/", getUser)       // GET /user/123
		r.Put("/", updateUser)    // PUT /user/123
		r.Delete("/", deleteUser) // DELETE /user/123
	})

	return r
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("All users here!")))
}

func createUser(w http.ResponseWriter, r *http.Request) {

	log.Printf("performed createUser")
	defer log.Printf("performed createUser END")

	var usr model.User

	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		log.Printf("Problem saving User ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usr.GeneratePassword()
	usr.Created = time.Now()
	usr.UsrCreated = "1"

	var v model.Valid = &usr

	res, err := driver.Create(v, r)
	if err != nil {
		log.Printf("Problem saving Statement ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usr.ID = res.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(usr); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("getUser is here")))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PUT User here"))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user here"))
}
