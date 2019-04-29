package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// StmntType - describes move for statement
type StmntType struct {
	ID         primitive.ObjectID `json:"_id,omitempty"`
	Name       string             `json:"name"`
	Comment    string             `json:"comment"`
	Code       string             `json:"code"`
	Investment bool               `json:"investment"`
	Loan       bool               `json:"loan"`
}
