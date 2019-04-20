package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/gmdmgithub/budget/config"
)

func main() {
	fmt.Println("Hi there use me  \u2318")

	cfg := config.Load()
	log.Printf("config read %+v", cfg)
}
