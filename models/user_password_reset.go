package models

import ()

//go:generate ffjson $GOFILE

// UserPasswordReset used for updating their password when they forgot it
type UserPasswordReset struct {
	TableName TableName `sql:"users,alias:user"       json:"-" lorem:"-"`

	Email       string `json:"email"  lorem:"email"`
	NewPassword string `json:"newPassword" lorem:"word,8,10"`
	Token       string `json:"token" lorem:"-"`
}
