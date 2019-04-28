package model

// Institution - model for institution like bank, Broker
type Institution struct {
	ID      string `json:"_id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Type    string `json:"type"`
	Audit
}
