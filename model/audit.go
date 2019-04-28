package model

import (
	"time"
)

// Audit - struct for audit data in db
type Audit struct {
	Created    time.Time `json:"created" bson:"created"`
	Updated    time.Time `json:"updated,omitempty" bson:"updated,omitempty"`
	UsrCreated string    `json:"usr_created" bson:"usr_created"`
	UsrUpdated string    `json:"usr_updated,omitempty" bson:"usr_updated,omitempty"`
}
