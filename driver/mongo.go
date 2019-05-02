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
func ConnectMgo(ctx context.Context, cfg *config.Config, db *DB) error {

	// mongodb://[username:password@]host[:port][/[database][?options]]
	//
	uri := fmt.Sprintf(
		"mongodb://%s:%s",
		cfg.DBS["MONGODB"].Host,
		cfg.DBS["MONGODB"].Port,
	)

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
	)

	db.C = ctx

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	log.Print("Connected to the DB!")
	db.Mongodb = client.Database(cfg.DBName)

	return nil
}

// Create - InsertOne, general method for model element to avoid repeated code
func Create(m model.Modeler) (*mongo.InsertOneResult, error) {

	if err := m.OK(); err != nil {
		log.Printf("Problem saving %T ... %+v", m, err)
		return nil, err
	}
	db := DBConn.Mongodb

	// store v in DB - next step
	res, err := db.Collection(m.ColName()).InsertOne(DBConn.C, m)
	if err != nil {
		log.Printf("Problem saving %T ... %+v", m, err)
		return nil, err
	}

	return res, nil
}

func DoOne(m model.Modeler, ID string, next func(m model.Modeler, ID string, filter bson.M) error) error {

	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Printf("Object %T with ID: %s and error: %v", m, ID, err.Error())
		return err
	}

	log.Printf("Getting object %T with ID: %v", m, _id)

	filter := bson.M{"_id": _id}

	return next(m, ID, filter)

}

// GetOne - gets one element of specific ID from collection named from struct
func GetOne(m model.Modeler, ID string, filter bson.M) error {

	db := DBConn.Mongodb
	err := db.Collection(m.ColName()).FindOne(DBConn.C, filter).Decode(m)
	if err != nil {
		log.Printf("Object %T with ID: %s and error: %v", m, ID, err.Error())
		return err
	}
	return nil
}

func DeleteOne(m model.Modeler, ID primitive.ObjectID) (*mongo.DeleteResult, error) {

	db := DBConn.Mongodb

	filter := bson.M{"_id": ID}
	// filter := bson.D{primitive.E{Key: "_id", Value: id}} //- another way to set filter with ID
	del, err := db.Collection(m.ColName()).DeleteOne(DBConn.C, filter)
	if err != nil {
		log.Printf("Object %T with ID: %s and error: %v", m, ID, err.Error())
		return nil, err
	}

	return del, nil
}

// UpdateOne - update document - all fields
func UpdateOne(m model.Modeler, ID primitive.ObjectID) (*mongo.UpdateResult, error) {

	db := DBConn.Mongodb

	filter := bson.M{"_id": ID}
	// update document
	update := bson.D{primitive.E{Key: "$set", Value: m}}
	res, err := db.Collection(m.ColName()).UpdateOne(DBConn.C, filter, update)
	if err != nil {
		log.Printf("Object %T with ID: %s and error: %v", m, ID, err.Error())
		return nil, err
	}

	return res, nil
}

func GetAllStatements() ([]model.Statement, error) {

	db := DBConn.Mongodb
	cursor, err := db.Collection("statements").Find(DBConn.C, bson.D{})
	if err != nil {
		log.Printf("GetAllStatements error: %+v", err.Error())
		return nil, err
	}
	defer cursor.Close(DBConn.C)

	// iterate through all documents
	var ms []model.Statement
	for cursor.Next(DBConn.C) {
		var m model.Statement
		// decode the document
		if err := cursor.Decode(&m); err != nil {
			log.Printf("GetAllStatements error: %+v", err.Error())
			return nil, err
		}
		// fmt.Printf("model: %+v", m)
		ms = append(ms, m)
	}
	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		log.Printf("GetAllStatements error: %+v", err.Error())
		return nil, err
	}

	return ms, nil
}

func GetAllCurrencies() ([]model.Currency, error) {

	db := DBConn.Mongodb
	// opts := options.Find()
	ctx := context.Background()
	cursor, err := db.Collection("currencies").Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Start and GetAllCurrencies error: %+v", err.Error())
		return nil, err
	}

	defer cursor.Close(ctx)

	// iterate through all documents
	var mc []model.Currency
	for cursor.Next(ctx) {
		var m model.Currency
		// decode the document
		if err := cursor.Decode(&m); err != nil {
			log.Printf("Problem with next - GetAllCurrencies error: %+v", err.Error())
			return nil, err
		}
		// fmt.Printf("model: %+v", m)
		mc = append(mc, m)
	}
	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		log.Printf("GetAllCurrencies error: %+v", err.Error())
		return nil, err
	}

	return mc, nil
}
