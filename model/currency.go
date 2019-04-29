package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Currency - represents exchange currency
type Currency struct {
	ID           primitive.ObjectID `json:"_id,omitempty"`
	Code         string             `json:"code"`
	ExchangeRate float64            `json:"exchange_rate"`
	Date         time.Time          `json:"date"`
	Base         bool               `json:"base"`
}
