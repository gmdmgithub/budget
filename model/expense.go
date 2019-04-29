package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Expense - represents expenses form the budget
type Expense struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	Name     string             `json:"name"`
	Comment  string             `json:"comment"`
	Period   string             `json:"period"`
	OneTime  bool               `json:"one_time"`
	Amount   int                `json:"amount"`
	Currency string             `json:"currency"`
	Date     time.Time          `json:"date"`
	Audit
}
