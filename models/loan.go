package models

import "time"

// Loan - represents a load object
type Loan struct {
	ID           string    `json:"_id"`
	Name         string    `json:"name"`
	Comment      string    `json:"comment"`
	Type         string    `json:"type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	InterestRate float64   `json:"interest_rate"`
	Amount       int       `json:"start_amount"`
	Currency     string    `json:"currency"`
	Active       bool      `json:"active"`
	Audit
}
