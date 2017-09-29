package models

// "github.com/axiomzen/null"

//go:generate ffjson $GOFILE

// FacebookUser is info that we get from facebook
type FacebookUser struct {
	FacebookID       string `sql:",null" form:"facebookId"       json:"facebookId"       lorem:"word,2,10"`
	FacebookUsername string `sql:",null" form:"facebookUsername" json:"facebookUsername" lorem:"word,2,10"`
	FacebookToken    string `sql:",null" form:"facebookToken"    json:"facebookToken"    lorem:"-"`
	FacebookEmail    string `sql:",null" form:"facebookEmail"    json:"facebookEmail"    lorem:"email"`
	FacebookPicture  string `sql:",null" form:"facebookPicture"    json:"facebookPicture"    lorem:"url"`
}
