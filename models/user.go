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
		Id:              user.ID,
		AuthToken:       user.AuthToken,
		Email:           user.Email,
		Verified:        user.Verified,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Status:          protobuf.UserStatus_created,
		FacebookID:      user.FacebookID,
		UserName:        user.UserName,
		FacebookPicture: user.FacebookPicture,
		FacebookToken:   user.FacebookToken,
		FacebookEmail: 	 user.FacebookEmail,
	}, nil
}
func (user *User) ProtobufPublic() (*protobuf.UserPublic, error) {
	return &protobuf.UserPublic{
		Id:              user.ID,
		Email:           user.Email,
		Status:          protobuf.UserStatus_created,
		FacebookID:      user.FacebookID,
		UserName:        user.UserName,
		FacebookPicture: user.FacebookPicture,
	}, nil
}
func (user *User) Merge(mergeWith *User) {
	// For linking email accounts
	if user.Email == "" {
		user.Email = mergeWith.Email
	}
	if user.Hash != nil && *user.Hash != "" {
		user.Hash = mergeWith.Hash
	}
	if user.VerifyEmailToken != "" {
		user.VerifyEmailToken = mergeWith.VerifyEmailToken
	}

	if user.UserName == "" {
		user.UserName = mergeWith.UserName
	}

	// For linking facebook accounts
	if user.FacebookID == "" {
		user.FacebookID = mergeWith.FacebookID
	}
	if user.FacebookUsername == "" {
		user.FacebookUsername = mergeWith.FacebookUsername
	}
	if user.FacebookToken == "" {
		user.FacebookToken = mergeWith.FacebookToken
	}
	if user.FacebookEmail == "" {
		user.FacebookEmail = mergeWith.FacebookEmail
	}
	if user.FacebookPicture == "" {
		user.FacebookPicture = mergeWith.FacebookPicture
	}

}

// Users is a slice of User pointers
// currently unused as we don't have any routes to paginate users yet
type Users []*User

func (users *Users) ProtobufPublic() (*protobuf.UsersPublic, error) {
	var protoUsers []*protobuf.UserPublic
	for _, user := range *users {
		protoUser, err := user.ProtobufPublic()
		if err != nil {
			return nil, err
		}
		protoUsers = append(protoUsers, protoUser)
	}
	return &protobuf.UsersPublic{
		Users: protoUsers,
	}, nil
}
