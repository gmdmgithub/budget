package model

import (
	"time"
)

// Audit - struct for audit data in db
type Audit struct {
	Created    time.Time `json:"created,omitempty" bson:"created,omitempty"`
	Updated    time.Time `json:"updated,omitempty" bson:"updated,omitempty"`
	UsrCreated string    `json:"usr_created,omitempty" bson:"usr_created,omitempty"`
	UsrUpdated string    `json:"usr_updated,omitempty" bson:"usr_updated,omitempty"`
}
