package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Institution - model for institution like bank, Broker
type Institution struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Code    string             `json:"code" bson:"code"`
	Name    string             `json:"name" bson:"name"`
	Comment string             `json:"comment,omitempty" bson:"comment,omitempty"`
	Type    string             `json:"type" bson:"type"`
	Audit
}

func (i *Institution) OK() error {

	if i.Code == "" || i.Name == "" || i.Type == "" {
		return errors.New("Fill in all required fields")
	}

	return nil
}

func (i *Institution) ColName() string {
	return "institutions"
}
