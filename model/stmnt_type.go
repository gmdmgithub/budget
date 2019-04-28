package model

// StmntType - describes move for statement
type StmntType struct {
	ID         string `json:"_id"`
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	Code       string `json:"code"`
	Investment bool   `json:"investment"`
	Loan       bool   `json:"loan"`
}
