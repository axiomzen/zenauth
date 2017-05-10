package models

import ()

//go:generate ffjson $GOFILE

// UserChangeEmail everything you need for changing email, nothing you don't
type UserChangeEmail struct {
	TableName TableName `sql:"users"       json:"-" lorem:"-"`
	ID        string    `bson:"id" json:"id" lorem:"-"`

	Email string `bson:"email"               json:"email"  lorem:"email"`
}
