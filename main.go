package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/rs/zerolog/log"

	"github.com/gmdmgithub/budget/config"
	"github.com/gmdmgithub/budget/driver"
)

func main() {
	fmt.Println("Hi there use me  \u2318")
	if err := run(); err != nil {
		log.Printf("Fatal problem during initialization: %v\n", err)
		os.Exit(1)
	}

}
func run() error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.Load()
	log.Printf("config read %+v", cfg)

	db, err := driver.ConnectMgo(cfg)
	if err != nil {
		log.Printf("No DB opened %v", err)
		os.Exit(-1)
	}
	log.Printf("DB %+v", db.Mongodb.Name())

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi there, home dir"))
	})
	log.Printf("Service is running on port %s", cfg.HTTPPort)
	return http.ListenAndServe(cfg.HTTPPort, r)

}
