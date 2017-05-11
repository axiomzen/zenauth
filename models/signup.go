package models

//go:generate ffjson $GOFILE

// Signup is the struct sent to us when signing up for the first time
type Signup struct {
	Email    string `form:"email"            json:"email" lorem:"email"`
	Password string `form:"password"         json:"password" lorem:"word,8,32"`
}
