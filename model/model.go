package model

// Valid - interface for checking corectness of model objects
type Valid interface {
	OK() error
	ColName() string
}
