package model

import (
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User - sctruct for user data
type User struct {
	// primitive.ObjectID
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Login       string             `json:"login" bson:"login"`
	Password    string             `json:"password," bson:"password"`
	OldPassword string             `json:"old_password,omitempty" bson:"old_password,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Active      bool               `json:"active" bson:"active"`
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

func (u *User) GeneratePassword(old bool) {

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("addUser: cannot create a password for the user: %v", err)
	}
	if old {
		u.OldPassword = u.Password
	}
	u.Password = fmt.Sprintf("%s", password)
}

// ComparePassword - compare password from db at User with plain password
func (u *User) ComparePassword(plainPwd []byte) bool {

	byteHash := []byte(u.Password) //password is hashed
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
