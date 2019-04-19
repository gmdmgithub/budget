package models

import "time"

// Expense - represents expenses form the budget
type Expense struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Period      string    `json:"period"`
	OneTime     bool      `json:"one_time"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Date        time.Time `json:"date"`
}
