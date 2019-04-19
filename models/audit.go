package models

import "time"

// Audit - struct for audit data in db
type Audit struct {
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	UsrCreated string    `json:"usr_created"`
	UsrUpdated string    `json:"usr_updated"`
}
