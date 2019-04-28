package models

import "time"

// Statement - model for db object that hold investments and loans
type Statement struct {
	ID           string    `json:"_id" bson:"_id"`
	Name         string    `json:"name" bson:"name"`
	Comment      string    `json:"comment" bson:"comment"`
	InstCode     string    `json:"inst_code" bson:"inst_code"`
	TypeCode     string    `json:"type_code" bson:"type_code"`
	StartDate    time.Time `json:"start_date" bson:"start_date"`
	EndDate      time.Time `json:"end_date" bson:"end_date"`
	InterestRate float64   `json:"interest_rate" bson:"interest_rate"`
	Amount       int       `json:"amount" bson:"amount"`
	CurrencyCode string    `json:"currency_code" bson:"currency_code"`
	Active       bool      `json:"active" bson:"active"`
	Audit
}
