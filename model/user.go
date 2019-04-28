package model

type User struct {
	// primitive.ObjectID
	ID       string `json:"_id" bson:"_id"`
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Age      int    `json:"age,omitempty" bson:"age,omitempty"`
	// CreatedAt primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
