package driver

import (
	"net/url"
	"time"

	"github.com/gmdmgithub/budget/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

// StatementsRange - gets statemts in range of data
func StatementsRange(r url.Values) ([]model.Statement, error) {
	filter := bson.M{}

	layout := "2006-01-02"

	df := true
	dt := true

	sf := r.Get("dateFrom")
	tf, err := time.Parse(layout, sf)
	if err != nil {
		df = false
		log.Printf("Problem with parsing %v", err)
	}

	st := r.Get("dateTo")
	tt, err := time.Parse(layout, st)
	if err != nil {
		dt = false
		log.Printf("Problem with parsing %v", err)
	}
	log.Printf("gt is %v and dt is %v", df, dt)

	if df && dt {
		filter["start_date"] = bson.M{"$gte": tf, "$lte": tt}
		// or := []bson.M{}
		// or = append(or, bson.M{"end_date": bson.M{"$gte": tt}})
		// or = append(or, bson.M{"end_date": nil})
		// filter["$or"] = or
	}

	if df && !dt {
		filter["start_date"] = bson.M{"$gte": tf}
	}
	// filter := bson.M{"name": bson.M{"$elemMatch": bson.M{"$eq": "golang"}}}
	if !df && dt {
		// or := []bson.M{}
		// or = append(or, bson.M{"end_date": bson.M{"$gte": tt}})
		// or = append(or, bson.M{"end_date": nil})
		// filter["$or"] = or

		filter["start_date"] = bson.M{"$lte": tt}
	}

	// log.Printf("values %v %v and filter %+v", tf, tt, filter)

	return filterStatements(filter)

}

// GetAllStatements - gets all statements in db
func GetAllStatements() ([]model.Statement, error) {
	return filterStatements(bson.M{})
}

// StatementsOnDate get statements from date
func StatementsOnDate(r url.Values) ([]model.Statement, error) {
	filter := bson.M{}

	layout := "2006-01-02"
	s := r.Get("date")
	t, err := time.Parse(layout, s)
	if err != nil {
		log.Printf("Problem with parsing %v", err)
		return nil, err
	}

	filter["start_date"] = bson.M{"$lte": t}
	or := []bson.M{}
	or = append(or, bson.M{"end_date": bson.M{"$gte": t}})
	or = append(or, bson.M{"end_date": nil})
	filter["$or"] = or

	return filterStatements(filter)
}

func filterStatements(filter bson.M) ([]model.Statement, error) {

	db := DBConn.Mongodb
	cursor, err := db.Collection("statements").Find(DBConn.C, filter)
	if err != nil {
		log.Printf("filterStatements error: %+v", err.Error())
		return nil, err
	}
	defer cursor.Close(DBConn.C)

	// iterate through all documents
	var ms []model.Statement
	for cursor.Next(DBConn.C) {
		var m model.Statement
		// decode the document
		if err := cursor.Decode(&m); err != nil {
			log.Printf("filterStatements error: %+v", err.Error())
			return nil, err
		}
		// fmt.Printf("model: %+v", m)
		ms = append(ms, m)
	}
	log.Printf("Res size %d", len(ms))
	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		log.Printf("filterStatements error: %+v", err.Error())
		return nil, err
	}

	return ms, nil
}
