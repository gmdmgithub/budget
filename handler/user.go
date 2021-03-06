package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gmdmgithub/budget/driver"
	"github.com/gmdmgithub/budget/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

const userContext = contextKey("user")

func UserRouter() http.Handler {

	r := chi.NewRouter()

	r.Get("/", getAllUsers)
	r.Post("/", createUser)
	r.Post("/login", login) // POST /login
	r.Route("/{usrID}", func(r chi.Router) {
		r.Use(usrContext)
		r.Get("/", getUser) // GET /user/123

		r.Put("/", updateUser)    // PUT /user/123
		r.Delete("/", deleteUser) // DELETE /user/123
	})

	return r
}

func usrContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usrID := chi.URLParam(r, "usrID")

		var usr model.User
		var v model.Modeler = &usr

		err := driver.DoOne(v, usrID, driver.GetOne)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), userContext, &usr)
		log.Printf("Data from DB: %+v with ID: %v", usr, usrID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	filter := bson.M{}
	options := options.Find()

	var u model.User

	res, err := driver.GetList(filter, options, &u)
	if err != nil {
		log.Printf("Problem to get the list of users %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var us []model.User
	for _, ui := range res {
		user, ok := ui.(*model.User)
		if ok {
			us = append(us, *user)
		}
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(us); err != nil {
		log.Printf("Problem to get the list of users %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte(fmt.Sprintf("All users here!")))
}

func login(w http.ResponseWriter, r *http.Request) {

	//another way to read data from body
	rlt, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Problem login User ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var u model.User

	err = json.Unmarshal(rlt, &u)
	if err != nil {
		log.Printf("Problem login User ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawPassword := u.Password

	//preferred way
	// if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
	// 	log.Printf("Problem saving User ... %v \n %+v\n", err, r.Body)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	log.Printf("User is: %+v", u)
	// filter posts tagged as golang
	// filter := bson.M{"login": bson.M{"$eq": u.Login}}
	filter := bson.M{}
	filter["login"] = u.Login

	// find one document

	db := driver.DBConn

	if err := db.Mongodb.Collection(u.ColName()).FindOne(db.C, filter).Decode(&u); err != nil {
		log.Printf("Problem login User ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("post: %+v\n", u)

	// check if password is correct
	if !u.ComparePassword([]byte(rawPassword)) {
		log.Printf("Problem login User %s - incorrect password: %s", u.Login, u.Password)
		http.Error(w, "Access dinied", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

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

	usr.GeneratePassword(false)
	usr.Created = time.Now()
	usr.UsrCreated = "1" //TODO change in the future

	var v model.Modeler = &usr

	res, err := driver.Create(v)
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

	ctx := r.Context()

	usr, ok := ctx.Value(userContext).(*model.User)
	if !ok {
		log.Printf("Problem with user context... ")
		http.Error(w, http.StatusText(422), 422)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(usr); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte(fmt.Sprintf("getUser is here")))
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	dbUsr, ok := ctx.Value(userContext).(*model.User)
	if !ok {
		log.Printf("Problem with user context... ")
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// take user from body - this is important

	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Problem with user decode... %v", err)
		http.Error(w, http.StatusText(422), 422)
		return
	}

	user.ID = dbUsr.ID
	user.Password = dbUsr.Password
	user.Created = dbUsr.Created
	user.UsrCreated = dbUsr.UsrCreated
	user.Updated = time.Now()
	user.UsrUpdated = "2" //correct in the future

	// update only with passward and login cannot be change
	user.Login = dbUsr.Login
	// first chcek if old passwod is the same as stored in DB
	if !dbUsr.ComparePassword([]byte(user.OldPassword)) {
		log.Printf("All password does not matches the old one ... %v ", r.Body)
		http.Error(w, "Old password problem", http.StatusNotAcceptable)
		return
	}
	// now plain password may be replaced
	user.GeneratePassword(true)

	var m model.Modeler = &user
	res, err := driver.UpdateOne(m, user.ID)

	if err != nil {
		log.Printf("Problem update User ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem update ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte("PUT User here"))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	usr, ok := ctx.Value(userContext).(*model.User)
	if !ok {
		log.Printf("Problem with user context... ")
		http.Error(w, http.StatusText(422), 422)
		return
	}
	var m model.Modeler = usr
	res, err := driver.DeleteOne(m, usr.ID)
	if err != nil {
		log.Printf(" delete Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	}

	// w.Write([]byte("Delete user here"))
}
