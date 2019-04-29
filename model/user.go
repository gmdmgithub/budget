package model

import "errors"

// User - sctruct for user data
type User struct {
	// primitive.ObjectID
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`
	Audit
}

// OK - check correcteness
func (u *User) OK() error {

	if u.Login == "" || u.Password == "" || u.Type == "" {
		return errors.New("Fill in all required fields")
	}

	return nil
}

// ColName - return name of collection in DB
func (u *User) ColName() string {

	return "users"
}
