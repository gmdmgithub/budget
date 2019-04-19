package models

// Institution - model for institution like bank, Broker
type Institution struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
