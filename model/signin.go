package model

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
