package driver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gmdmgithub/budget/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//DB struct of databases in project
type DB struct {
	Mongodb *mongo.Database
	SQL     *sql.DB //if any
}

// DBConn hold connection to the databases
var dbConn = &DB{}

// ConnectSQL - connect to mySQL DB
func ConnectSQL(host, port, user, pass, name string) (*DB, error) {

	//  username:password@protocol(address)/dbname?param=value
	// https://github.com/go-sql-driver/mysql
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		user,
		pass,
		host,
		port,
		name,
	)
	d, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err)
	}
	dbConn.SQL = d
	return dbConn, err
}

// ConnectMgo return connection to the mongodb
func ConnectMgo(cfg *config.Config) (*DB, error) {
	// mongodb://[username:password@]host[:port][/[database][?options]]
	//
	uri := fmt.Sprintf(
		"mongodb://%s:%s",
		cfg.DBS["MONGODB"].Host,
		cfg.DBS["MONGODB"].Port,
	)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Print("Connected to the DB!")
	dbConn.Mongodb = client.Database(cfg.DBName)

	return dbConn, err
}
