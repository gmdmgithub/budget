package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Institution - model for institution like bank, Broker
type Institution struct {
	ID      primitive.ObjectID `json:"_id,omitempty"`
	Code    string             `json:"code"`
	Name    string             `json:"name"`
	Comment string             `json:"comment"`
	Type    string             `json:"type"`
	Audit
}
