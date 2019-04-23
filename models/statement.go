package models

import "time"

// Statement - model for db object that hold investments and loans
type Statement struct {
	ID           string    `json:"_id"`
	Name         string    `json:"name"`
	Comment      string    `json:"comment"`
	InstCode     string    `json:"inst_code"`
	TypeCode     string    `json:"type_code"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	InterestRate float64   `json:"interest_rate"`
	Amount       int       `json:"amount"`
	CurrencyCode string    `json:"currency_code"`
	Active       bool      `json:"active"`
	Audit
}
