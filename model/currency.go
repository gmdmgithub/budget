package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Currency - represents exchange currency
type Currency struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Code         string             `bson:"code,omitempty"`
	ExchangeRate float64            `bson:"exchange_rate,omitempty"`
	Date         time.Time          `bson:"date,omitempty"`
	Base         bool               `bson:"base,omitempty"`
}

func (c *Currency) OK() error {
	if c.Code == "" || c.ExchangeRate <= 0 {
		return errors.New("Fill in required fields")
	}

	return nil
}

func (c *Currency) ColName() string {
	return "currencies"
}
