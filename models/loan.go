package models

import "time"

// Loan - represents a load object
type Loan struct {
	ID           string    `json:"_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         string    `json:"type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	InterestRate float64   `json:"interest_rate"`
	StartAmount  float64   `json:"start_amount"`
	EndAmount    float64   `json:"end_amount"`
	Currency     string    `json:"currency"`
	Active       bool      `json:"active"`
	Audit        Audit     `json:"audit"`
}
