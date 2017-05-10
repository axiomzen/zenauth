package models

//go:generate ffjson $GOFILE

// UserAuth holds the password as well
type UserAuth struct {
	UserBase
	Password string `json:"password" lorem:"word,8,10"`
}
