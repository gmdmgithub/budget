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

const statementContext = contextKey("statement")

// StatementRouter - a completely separate router for administrator routes
func StatementRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", allStatements)
	r.Post("/", createStatement())      // POST /articles - different way (func is returned)
	r.Get("/range", rangeStatements)    // GET /statement/data from date to
	r.Get("/on-date", onDateStatements) // GET /statement/data from date to
	r.Route("/{stID}", func(r chi.Router) {
		r.Use(statementCtx)
		r.Get("/", getStatement)       // GET /statement/123
		r.Put("/", updateStatement)    // PUT /statement/123
		r.Delete("/", deleteStatement) // DELETE /statement/123
		// Regexp url parameters:
		r.Get("/{name:[a-z-]+}", allStatements) // GET /statement/income from ABC

	})

	return r
}

// statementCtx add Statement to the context - or any other necessary objects
func statementCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stID := chi.URLParam(r, "stID")
		var stmt model.Statement
		err := driver.DoOne(&stmt, stID, driver.GetOne)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), statementContext, &stmt)
		log.Printf("Data from DB: %+v with ID: %v", stmt, stID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createStatement() http.HandlerFunc {
	// DoOnce part
	log.Printf("Hi there post func prepared")
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("performed create statement")
		defer log.Printf("performed create statement END")

		var stmt model.Statement

		if err := json.NewDecoder(r.Body).Decode(&stmt); err != nil {
			log.Printf("Problem saving User ... %v \n %+v\n", err, r.Body)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt.UsrCreated = "123456" //TODO temporary 1, should be current user
		stmt.Created = time.Now()

		res, err := driver.Create(&stmt)
		if err != nil {
			log.Printf("Problem saving Statement ... %v \n %+v\n", err, r.Body)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt.ID = res.InsertedID.(primitive.ObjectID)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stmt); err != nil {
			log.Printf(" json Problem ... %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func rangeStatements(w http.ResponseWriter, r *http.Request) {

	res, err := driver.StatementsRange(r.URL.Query())
	if err != nil {
		log.Printf("Problem with getting range res %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func onDateStatements(w http.ResponseWriter, r *http.Request) {
	res, err := driver.StatementsOnDate(r.URL.Query())
	if err != nil {
		log.Printf("Problem with getting range res %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func allStatements(w http.ResponseWriter, r *http.Request) {

	stmts, err := driver.GetAllStatements()
	if err != nil {
		log.Printf("allStatements problem ... %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stmts); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte(fmt.Sprintf("All statements here!")))
}

func getStatement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// check the type
	// log.Printf("what type? %T and values %+v", ctx.Value("statement"), ctx.Value("statement"))
	statement, ok := ctx.Value(statementContext).(*model.Statement)
	// statement, ok := ctx.Value("stID").(string)
	if !ok {
		log.Print("problem with statement")
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(statement); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func updateStatement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// check the type
	statement, ok := ctx.Value(statementContext).(*model.Statement)
	if !ok {
		log.Print("problem with statement")
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// first we read from DB but finally should be this from body
	var stmt model.Statement
	if err := json.NewDecoder(r.Body).Decode(&stmt); err != nil {
		log.Printf("Problem saving User ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stmt.ID = statement.ID
	stmt.Created = statement.Created
	stmt.UsrCreated = statement.UsrCreated
	stmt.Updated = time.Now()
	stmt.UsrUpdated = "1" //correct in the future

	res, err := driver.UpdateOne(&stmt, stmt.ID)
	if err != nil {
		log.Printf("Problem update Statement ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem update ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// w.Write([]byte(fmt.Sprintf("Update statement?")))
}
func deleteStatement(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// check the type
	// log.Printf("what type? %T and values %+v", ctx.Value("statement"), ctx.Value("statement"))
	statement, ok := ctx.Value(statementContext).(*model.Statement)
	// statement, ok := ctx.Value("stID").(string)
	if !ok {
		log.Print("problem with statement")
		http.Error(w, http.StatusText(422), 422)
		return
	}

	res, err := driver.DeleteOne(statement, statement.ID)
	if err != nil {
		log.Printf("Problem delete Statement ... %v \n %+v\n", err, r.Body)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(" json Problem ... %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// w.Write([]byte(fmt.Sprintf("Delete statement")))
}
