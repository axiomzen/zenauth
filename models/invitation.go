package models

import (
	"fmt"

	"github.com/axiomzen/null"
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/protobuf"
)

type InvitationRequest struct {
	InviteCodes []string `json:"inviteCodes"`
}
type InvitationResponse struct {
	Users []*protobuf.UserPublic `json:"users"`
}

type Invitation struct {
	ID        string    `sql:",pk"`
	TableName TableName `sql:"invitations,alias:invitation"`
	Type      string
	Code      string
	CreatedAt null.Time `sql:",null"`
}

type Invitations []*Invitation

func (invitation *Invitation) UserPublicProtobuf() (*protobuf.UserPublic, error) {
	user := &protobuf.UserPublic{
		Id:     invitation.ID,
		Status: protobuf.UserStatus_invited,
	}
	switch invitation.Type {
	case constants.InvitationTypeEmail:
		user.Email = invitation.Code
	case constants.InvitationTypeFacebook:
		user.FacebookID = invitation.Code
	}
	return user, nil
}

func (invitation *Invitation) UpdateUserWithInvitationInfo(user *User) error {
	switch invitation.Type {
	case constants.InvitationTypeEmail:
		user.Email = invitation.Code
	case constants.InvitationTypeFacebook:
		user.FacebookID = invitation.Code
	default:
		return fmt.Errorf("Invitation type %s not supported", invitation.Type)
	}
	return nil
}
