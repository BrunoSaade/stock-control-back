package model

type User struct {
	ID       string    `json:"id,omitempty" bson:"id"`
	Email    string    `json:"email,omitempty" bson:"email"`
	Password string    `json:"password,omitempty" bson:"password"`
	Stock    []Product `json:"stock,omitempty" bson:"stock,omitempty"`
}
