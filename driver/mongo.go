package driver

import (
	"context"
	"fmt"

	"github.com/gmdmgithub/budget/model"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gmdmgithub/budget/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectMgo - return connection to the mongodb
func ConnectMgo(ctx context.Context, cfg *config.Config) (*DB, error) {

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

// Create - InsertOne, general method for model element to avoid repeated code
func Create(v model.Valid) (*mongo.InsertOneResult, error) {

	if err := v.OK(); err != nil {
		log.Printf("Problem saving  %T ... %+v", v, err)
		return nil, err
	}
	db := DBConn.Mongodb

	// store v in DB - next step
	res, err := db.Collection(v.ColName()).InsertOne(DBConn.C, v)
	if err != nil {
		log.Printf("Problem saving %T ... %+v", v, err)
		return nil, err
	}

	return res, nil
}

// GetOne - gets one element of specific ID from collection named from struct
func GetOne(v model.Valid, ID string) error {

	db := DBConn.Mongodb

	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Printf("Object %T with ID: %s and error: %v", v, ID, err.Error())
		return err
	}

	log.Printf("Getting object %T with ID: %v", v, _id)

	filter := bson.M{"_id": _id}
	err = db.Collection(v.ColName()).FindOne(DBConn.C, filter).Decode(v)
	if err != nil {
		log.Printf("Object %T with ID: %s and error: %v", v, ID, err.Error())
		return err
	}
	return nil
}
