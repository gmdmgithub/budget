package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Statement - model for db object that hold investments and loans
type Statement struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Comment      string             `json:"comment" bson:"comment"`
	InstCode     string             `json:"inst_code" bson:"inst_code"`
	TypeCode     string             `json:"type_code" bson:"type_code"`
	StartDate    time.Time          `json:"start_date" bson:"start_date"`
	EndDate      time.Time          `json:"end_date,omitempty" bson:"end_date,omitempty"`
	InterestRate float64            `json:"interest_rate" bson:"interest_rate"`
	Amount       int                `json:"amount" bson:"amount"`
	CurrencyCode string             `json:"currency_code" bson:"currency_code"`
	Active       bool               `json:"active" bson:"active"`
	Audit
}

// OK - check correcteness
func (s *Statement) OK() error {
	if s.Name == "" || s.InstCode == "" || s.TypeCode == "" {
		return errors.New("Fill in all required fields")
	}

	return nil
}
