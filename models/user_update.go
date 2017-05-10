package models

import ()

//go:generate ffjson $GOFILE

// UserUpdate is all the fields that we update on a normal update (PUT)
type UserUpdate struct {
	TableName TableName `sql:"users"       json:"-" lorem:"-"`
	ID        string    `json:"id" lorem:"-"`

	Preferences map[string]string `bson:"-"          json:"preferences" lorem:"-"`
	FirstName   *string           `bson:"first_name"          json:"firstName"  lorem:"word,2,10"`
	LastName    *string           `bson:"last_name"           json:"lastName"  lorem:"word,2,10"`
}
