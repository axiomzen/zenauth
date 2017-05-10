package models

//go:generate ffjson $GOFILE

// UserChangePassword used for changing their password
// incoming from the request
type UserChangePassword struct {
	UserBase
	Hash        string `bson:"hash"                    json:"-"  lorem:"-"`
	NewPassword string `json:"newPassword" lorem:"word,8,10"`
	OldPassword string `json:"oldPassword" lorem:"word,8,10"`
}
