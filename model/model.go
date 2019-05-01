package model

// Modeler - interface for gathering of model objects
type Modeler interface {
	OK() error
	ColName() string
}
