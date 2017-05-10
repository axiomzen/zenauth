package models

import (
// "github.com/axiomzen/null"
)

//go:generate ffjson $GOFILE

// User struct holds our complete user information
type User struct {
	UserBase    `bson:",inline"`
	Preferences map[string]string `bson:"-" json:"preferences" lorem:"-"`
	ResetToken  *string           `bson:"reset_token" json:"-" lorem:"-"`
	Hash        *string           `bson:"hash" json:"-"  lorem:"-"`
	AuthToken   string            `bson:"-" json:"authToken,omitempty" lorem:"-" sql:"-"`

	// TODO: figure this out
	SocialLogin bool `bson:"social_login" json:"socialLogin"`
}

// Users is a slice of User pointers
// currently unused as we don't have any routes to paginate users yet
type Users []*User
