package models

import (
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

func (invitation *Invitation) UserPublicProtobuf() (*protobuf.UserPublic, error) {
	user := &protobuf.UserPublic{
		Id:     invitation.ID,
		Status: protobuf.UserStatus_invited,
	}
	if invitation.Type == constants.InvitationTypeEmail {
		user.Email = invitation.Code
	}
	return user, nil
}
