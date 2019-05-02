package driver

import (
	"database/sql"
	"fmt"
)

// ConnectSQL - connect to mySQL DB
func ConnectSQL(host, port, user, pass, name string, db *DB) error {

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
	var err error
	db.SQL, err = sql.Open("mysql", dbSource)
	if err != nil {
		panic(err)
	}
	return err
}
