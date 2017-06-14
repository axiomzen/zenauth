package v1

import (
	"strings"

	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/models"
	"github.com/axiomzen/zenauth/protobuf"
	"github.com/gocraft/web"
)

// InvitationContext for user authenticated routes
type InvitationContext struct {
	*UserContext
}

// Create invitation route creates multiple invitations
//
//   POST /users/invite/email
//
// Returns
//   201 Status Created
func (c *InvitationContext) Create(rw web.ResponseWriter, req *web.Request) {
	var invitationRequest models.InvitationRequest
	// decode request
	if !c.DecodeHelper(&invitationRequest, "Couldn't decode the request body", rw, req) {
		return
	}

	invitations := make([]*models.Invitation, len(invitationRequest.InviteCodes))
	var user models.User
	for idx, email := range invitationRequest.InviteCodes {
		// Verify invite email is valid
		if strings.Count(email, "@") != 1 {
			model := models.NewErrorResponse(constants.APIValidationEmailNotValid,
				models.NewAZError("Invalid email address"), "Could not create invitation")
			c.Render(constants.StatusBadRequest, model, rw, req)
			return
		}
		invitations[idx] = &models.Invitation{
			Type: constants.InvitationTypeEmail,
			Code: helpers.EmailSanitize(email),
		}
		// Verify we don't already have a user with this email
		user.Email = invitations[idx].Code

		// TODO: Should we error here? Or continue inviting the rest of the list?
		if err := c.DAL.GetUserByEmail(&user); err == nil {
			model := models.NewErrorResponse(constants.APIDatabaseCreateInvitation,
				models.NewAZError("User with email already exists"), "Could not create invitation")
			c.Render(constants.StatusBadRequest, model, rw, req)
			return
		}
	}

	if err := c.DAL.CreateInvitations(&invitations); err != nil {
		model := models.NewErrorResponse(constants.APIInvitationsCreationError, models.NewAZError(err.Error()), "unable to create the invitations")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	invitationResponse := models.InvitationResponse{
		Users: make([]*protobuf.UserPublic, len(invitations)),
	}
	var err error
	for idx, invitation := range invitations {
		invitationResponse.Users[idx], err = invitation.UserPublicProtobuf()
		if err != nil {
			model := models.NewErrorResponse(constants.APIInvitationsCreationError, models.NewAZError(err.Error()), "unable to get the view of the invitation")
			c.Render(constants.StatusInternalServerError, model, rw, req)
			return
		}
	}

	rw.Header().Set("Location", "/v1/users")

	c.Render(constants.StatusCreated, &invitationResponse, rw, req)
}
