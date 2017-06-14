package models

//go:generate ffjson $GOFILE

// FacebookSignup is the struct sent to us when signing up with Facebook
type FacebookSignup struct {
	FacebookUser
	//Link      bool    `form:"link"                 json:"link"`
	Email string `form:"email"                json:"email" lorem:"email"`
}
