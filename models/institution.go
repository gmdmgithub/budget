package models

// Institution - model for institution like bank, Broker
type Institution struct {
	ID      string `json:"_id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Type    string `json:"type"`
	Audit
}
