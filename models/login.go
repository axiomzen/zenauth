package models

//go:generate ffjson $GOFILE

// Login is the basic message sent to us upon logging into the app
type Login struct {
	Email    string `form:"email"        json:"email"       lorem:"email"`
	Password string `form:"password"     json:"password"    lorem:"word,2,10"`
	UserName string `form:"userName" json:"userName" lorem:"uuid"`
}
