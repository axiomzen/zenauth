package models

import ()

//go:generate ffjson $GOFILE

// Signup is the struct sent to us when signing up for the first time
type Signup struct {
	FirstName   *string `form:"firstName"        json:"firstName" lorem:"sentence,1,2"`
	LastName    *string `form:"lastName"         json:"lastName" lorem:"sentence,1,2"`
	Email       string  `form:"email"            json:"email" lorem:"email"`
	Password    string  `form:"password"         json:"password" lorem:"word,8,32"`
	Username    *string `form:"username"         json:"username" lorem:"word,2,10"`
	Description *string `form:"description"      json:"description" lorem:"paragraph,1,2"`
	Image       *string `form:"image"            json:"image" lorem:"url"`
	// TODO: needed?
	//TwitterHandle null.String `form:"twitterHandle"    json:"twitterHandle"`
}
