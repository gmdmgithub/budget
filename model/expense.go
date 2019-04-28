package model

import "time"

// Expense - represents expenses form the budget
type Expense struct {
	ID       string    `json:"_id"`
	Name     string    `json:"name"`
	Comment  string    `json:"comment"`
	Period   string    `json:"period"`
	OneTime  bool      `json:"one_time"`
	Amount   int       `json:"amount"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
	Audit
}
