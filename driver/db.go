package driver

import (
	"context"
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

//DB struct of databases in project
type DB struct {
	Mongodb *mongo.Database
	SQL     *sql.DB //if any
	C       context.Context
}

// DBConn hold connection to the databases
var DBConn = &DB{}
