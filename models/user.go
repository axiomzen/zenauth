package models

import (
	"github.com/axiomzen/zenauth/protobuf"
	gpPtypes "github.com/golang/protobuf/ptypes"
)

// "github.com/axiomzen/null"

//go:generate ffjson $GOFILE

// User struct holds our complete user information
type User struct {
	UserBase
	ResetToken       *string `json:"-" lorem:"-"`
	Hash             *string `json:"-" lorem:"-"`
	AuthToken        string  `json:"authToken,omitempty" lorem:"-" sql:"-"`
	VerifyEmailToken string  `json:"-" lorem:"-" sql:"-"`

	FacebookUser
}

func (user *User) Protobuf() (*protobuf.User, error) {
	createdAt, err := gpPtypes.TimestampProto(user.CreatedAt.Time)
	if err != nil {
		return nil, err
	}
	updatedAt, err := gpPtypes.TimestampProto(user.UpdatedAt.Time)
	if err != nil {
		return nil, err
	}
	return &protobuf.User{
		Id:         user.ID,
		AuthToken:  user.AuthToken,
		Email:      user.Email,
		Verified:   user.Verified,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Status:     protobuf.UserStatus_created,
		FacebookID: user.FacebookID,
	}, nil
}
func (user *User) ProtobufPublic() (*protobuf.UserPublic, error) {
	return &protobuf.UserPublic{
		Id:         user.ID,
		Email:      user.Email,
		Status:     protobuf.UserStatus_created,
		FacebookID: user.FacebookID,
	}, nil
}

// Users is a slice of User pointers
// currently unused as we don't have any routes to paginate users yet
type Users []*User
