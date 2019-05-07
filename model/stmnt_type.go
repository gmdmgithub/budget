package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StmntType - describes move for statement
type StmntType struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Comment    string             `json:"comment,omitempty" bson:"comment,omitempty"`
	Code       string             `json:"code" bson:"code"`
	Investment bool               `json:"investment" bson:"investment"`
	Loan       bool               `json:"loan" bson:"loan"`
	Audit
}

func (s *StmntType) OK() error {
	if s.Name == "" || s.Code == "" {
		return errors.New("Fill in all required fields")
	}

	return nil
}

func (s *StmntType) ColName() string {
	return "stmtTypes"
}
