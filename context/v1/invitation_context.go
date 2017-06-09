package v1

import (
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/models"
	"github.com/axiomzen/zenauth/protobuf"
	"github.com/gocraft/web"
)

// InvitationContext for user authenticated routes
type InvitationContext struct {
	*UserContext
}

// Create invitation route
//
//   POST /users/invite
//
// Returns
//   201 Status Created
func (c *UserContext) Create(rw web.ResponseWriter, req *web.Request) {
	var invitationRequest models.InvitationRequest
	// decode request
	if !c.DecodeHelper(&invitationRequest, "Couldn't decode the request body", rw, req) {
		return
	}

	invitations := make([]*models.Invitation, len(invitationRequest.Emails))
	for idx, email := range invitationRequest.Emails {
		invitations[idx] = &models.Invitation{
			Email: email,
		}
	}

	if err := c.DAL.GetOrCreateInvitations(&invitations); err != nil {
		model := models.NewErrorResponse(constants.APIInvitationsCreationError, models.NewAZError(err.Error()), "unable to create the invitations")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	invitationResponse := models.InvitationResponse{
		Users: make([]*protobuf.UserPublic, len(invitations)),
	}
	for idx, invitation := range invitations {
		invitationResponse.Users[idx] = invitation.UserPublicProtobuf()
	}

	rw.Header().Set("Location", "/v1/users")

	c.Render(constants.StatusCreated, &invitationResponse, rw, req)
}
