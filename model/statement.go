package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Statement - model for db object that hold investments and loans
type Statement struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Comment      string             `bson:"comment"`
	InstCode     string             `bson:"inst_code"`
	TypeCode     string             `bson:"type_code"`
	StartDate    time.Time          `bson:"start_date"`
	EndDate      time.Time          `bson:"end_date,omitempty"`
	InterestRate float64            `bson:"interest_rate"`
	Amount       int                `bson:"amount"`
	CurrencyCode string             `bson:"currency_code"`
	Active       bool               `bson:"active"`
	Audit
}

// OK - check correcteness
func (s *Statement) OK() error {
	if s.Name == "" || s.InstCode == "" || s.TypeCode == "" {
		return errors.New("Fill in all required fields")
	}

	return nil
}

// ColName - return name of collection in DB
func (s *Statement) ColName() string {
	return "statements"
}
