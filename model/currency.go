package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Currency - represents exchange currency
type Currency struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Code         string             `json:"code" bson:"code"`
	ExchangeRate float64            `json:"exchange_rate" bson:"exchange_rate"`
	Date         time.Time          `json:"date" bson:"date"`
	Base         bool               `json:"base" bson:"base"`
}

func (c *Currency) OK() error {
	if c.Code == "" || c.ExchangeRate <= 0 {
		return errors.New("Fill in required fields")
	}

	return nil
}

func (c *Currency) ColName() string{
	return "currencies"
}

func 
