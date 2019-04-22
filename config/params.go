package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const ()

// DB - struct describing DBs
type DB struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Config - struct with coonfiguration
type Config struct {
	AccessToken string
	HTTPPort    string
	DBName      string
	DBS         map[string]DB
}

// Load - gather config data
func Load() (c *Config) {

	// first read .env file and put it to env
	if err := godotenv.Load(); err != nil {
		log.Printf("Fatal problem during initialization: %v\n", err)
		os.Exit(1)
	}
	// set logger configuarion
	LoadLog()

	c = &Config{}
	p, ok := os.LookupEnv("HTTP_PORT")
	if !ok {
		log.Print("No http port in .env file, default 8000 taken")
		p = ":8000"
	}
	c.HTTPPort = p
	dn, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Print("No db name in .env file, default \"budget\" taken")
		dn = "budget"
	}
	c.DBName = dn
	c.AccessToken = os.Getenv("ACCESS_TOKEN") //no default value for token

	dbm := make(map[string]DB)
	// set config to all DBS
	dbm["MONGODB"] = config("MONGODB")
	dbm["MYSQL"] = config("MYSQL")

	c.DBS = dbm

	return c
}

func config(name string) DB {

	var ok bool
	var d DB
	d.Host, ok = os.LookupEnv(name + "_HOST")
	if !ok {
		log.Print("No DB host in .env file aborted")
		return d
	}
	d.Port, ok = os.LookupEnv(name + "_PORT")
	if !ok {
		log.Print("No DB port in .env file aborted")
		return d
	}
	d.User = os.Getenv(name + "_USER")
	d.Password = os.Getenv(name + "_PASSWORD")

	return d
}
