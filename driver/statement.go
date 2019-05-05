package driver

import (
	"net/url"
	"time"

	"github.com/gmdmgithub/budget/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StatementsRange - gets statements in range of data
func StatementsRange(r url.Values) ([]model.Statement, error) {
	filter := bson.M{}
	options := options.Find()

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
	// log.Printf("gt is %v and dt is %v", df, dt)

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
		filter["start_date"] = bson.M{"$lte": tt}
	}
	return filterStatements(filter, options)
}

// GetAllStatements - gets all statements in db, no filter no options
func GetAllStatements() ([]model.Statement, error) {

	return filterStatements(bson.M{}, options.Find())
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

	options := options.Find()

	return filterStatements(filter, options)
}

func filterStatements(filter bson.M, options *options.FindOptions) ([]model.Statement, error) {

	var s model.Statement
	res, err := GetList(filter, options, &s)
	if err != nil {
		log.Printf("Problem with get statements %v", res)
		return nil, err
	}
	var ss []model.Statement
	for _, sI := range res {
		ss = append(ss, *sI.(*model.Statement))
	}
	return ss, nil
}
