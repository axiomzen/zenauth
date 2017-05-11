package models

// "github.com/axiomzen/null"

//go:generate ffjson $GOFILE

// User struct holds our complete user information
type User struct {
	UserBase
	ResetToken *string `json:"-" lorem:"-"`
	Hash       *string `json:"-" lorem:"-"`
	AuthToken  string  `json:"authToken,omitempty" lorem:"-" sql:"-"`
}

// Users is a slice of User pointers
// currently unused as we don't have any routes to paginate users yet
type Users []*User
