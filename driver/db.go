package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gmdmgithub/budget/model"

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
	C       context.Context
}

// DBConn hold connection to the databases
var DBConn = &DB{}

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
	DBConn.SQL = d
	return DBConn, err
}

// ConnectMgo return connection to the mongodb
func ConnectMgo(cfg *config.Config, ctx context.Context) (*DB, error) {
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

	DBConn.C = ctx

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Print("Connected to the DB!")
	DBConn.Mongodb = client.Database(cfg.DBName)

	return DBConn, err
}

func Create(v model.Valid, r *http.Request) (*mongo.InsertOneResult, error) {

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Problem saving Statement ... %v \n %+v\n", err, r.Body)
		return nil, err
	}

	if err := v.OK(); err != nil {
		log.Printf("Problem saving Statement ... %v \n %+v\n", err, r.Body)
		return nil, err
	}
	db := DBConn.Mongodb

	// store v in DB - next step
	log.Printf("DB %v", DBConn.Mongodb.Name())
	res, err := db.Collection(v.ColName()).InsertOne(DBConn.C, v)

	if err != nil {
		log.Printf("Problem saving Statement ... %v \n %+v\n", err, r.Body)
		return nil, err
	}

	return res, nil
}
