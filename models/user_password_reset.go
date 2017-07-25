package models

//go:generate ffjson $GOFILE

// UserPasswordReset used for updating their password when they forgot it
type UserPasswordReset struct {
	TableName TableName `sql:"users,alias:user"       json:"-" lorem:"-"`

	Email       string `json:"email" form:"email" lorem:"email"`
	NewPassword string `json:"newPassword" form:"newPassword" lorem:"word,8,10"`
	Token       string `json:"token" form:"token" lorem:"-"`
	Redirect    string `json:"redirect" form:"redirect" lorem:"-"`
}
