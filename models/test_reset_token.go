package models

//go:generate ffjson $GOFILE

// TestResetToken is the response for the test route to get the users reset token
type TestResetToken struct {
	Token string `form:"token"   json:"token"`
}
