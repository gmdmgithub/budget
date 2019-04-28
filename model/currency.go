package model

import "time"

// Currency - represents exchange currency
type Currency struct {
	ID           string    `json:"_id"`
	Code         string    `json:"code"`
	ExchangeRate float64   `json:"exchange_rate"`
	Date         time.Time `json:"date"`
	Base         bool      `json:"base"`
}
