package models

import "github.com/axiomzen/zenauth/protobuf"
import "github.com/axiomzen/null"

type InvitationRequest struct {
	Emails []string `json:"emails"`
}
type InvitationResponse struct {
	Users []*protobuf.UserPublic `json:"users"`
}

type Invitation struct {
	ID        string
	Email     string
	CreatedAt null.Time `sql:",null"`
}

func (invitation *Invitation) UserPublicProtobuf() *protobuf.UserPublic {
	return &protobuf.UserPublic{
		Id:     invitation.ID,
		Email:  invitation.Email,
		Status: protobuf.UserStatus_invited,
	}
}
