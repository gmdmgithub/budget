package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/gmdmgithub/budget/config"
	"github.com/gmdmgithub/budget/driver"
)

func main() {
	fmt.Println("Hi there use me  \u2318")

	cfg := config.Load()
	log.Printf("config read %+v", cfg)

	db, err := driver.ConnectMgo(cfg)
	if err != nil {
		log.Printf("No DB opened %v", err)
		os.Exit(-1)
	}
	log.Printf("DB %+v", db.Mongodb.Name())
}
