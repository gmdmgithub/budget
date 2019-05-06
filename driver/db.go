package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

//DB struct of databases in project
type DB struct {
	Mongodb *mongo.Database
	SQL     *sql.DB //if any
	C       context.Context
}

// DBConn hold connection to the databases
var DBConn *DB

// deepCopy is essential for copying data between two interfaces (as a data copy)
// alternalively copy values of struct - but this is general
func deepCopy(v interface{}) (interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	r := reflect.New(reflect.TypeOf(v))
	err = json.Unmarshal(data, r.Interface())
	if err != nil {
		return nil, err
	}
	return r.Elem().Interface(), err
}
