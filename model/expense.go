package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Expense - represents expenses form the budget
type Expense struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Comment  string             `json:"comment,omitempty" bson:"comment,omitempty"`
	Period   string             `json:"period" bson:"period"`
	OneTime  bool               `json:"one_time" bson:"one_time"`
	Amount   int                `json:"amount" bson:"amount"`
	Currency string             `json:"currency" bson:"currency"`
	Date     time.Time          `json:"date" bson:"date"`
	Audit
}

func (e *Expense) OK() error {

	if e.Name == "" || e.Amount <= 0 || e.Currency == "" {
		return errors.New("Fill in required data")
	}
	return nil
}

func (e *Expense) ColName() string {
	return "expenses"
}
