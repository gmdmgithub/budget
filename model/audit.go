package model

import (
	"time"
)

// Audit - struct for audit data in db
type Audit struct {
	Created    time.Time `bson:"created,omitempty"`
	Updated    time.Time `bson:"updated,omitempty"`
	UsrCreated string    `bson:"usr_created,omitempty"`
	UsrUpdated string    `bson:"usr_updated,omitempty"`
}
