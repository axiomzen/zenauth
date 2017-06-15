package v1

import (
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/data"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/models"
	"github.com/gocraft/web"

	"fmt"
)

// TestContext for testing routes
type TestContext struct {
	*APIAuthContext
}

// Panic will cause a panic on purpose
//
//   GET /panic
//
// Returns
//   500 error
func (c *TestContext) Panic(rw web.ResponseWriter, req *web.Request) {
	// not that it matters, but 99% of the time its a nil pointer dereference
	var thisWillPanic *bool
	*thisWillPanic = true
}

// TODO: decide if we want this or not
// // UserAuthRequired Middleware: Authorizes a user by authenticating the Json Web Token
// func (c *TestContext) UserAuthRequired(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

//     jwt := helpers.JWTHelper{c.Config.HashSecretBytes, r.Header.Get(c.Config.AuthTokenHeader)}

// 	jwtTokenResult := jwt.Verify(c.Config.JwtClaimUserID)

// 	switch jwtTokenResult.Status {
// 	case helpers.JWTokenStatusValid:
// 		if !c.DAL.IsValidId(userId) {
// 			c.UnauthorizedHandler(w, r)
// 			return
// 		} else {
// 			c.UserId = userId
// 			c.Log = c.Log.WithField("userId", userId)
// 		}
// 		next(w, r)
// 	case helpers.JWTokenStatusExpired:
// 		c.ExpiredHandler(w, r)
// 	case helpers.JWTokenStatusInvalid:
// 		c.UnauthorizedHandler(w, r)
// 	}
// }

// UserPasswordResetTokenGet gets a users password reset token
//
//   GET /test/users/ResourcePasswordReset
//
// Returns
//   200 OK
func (c *TestContext) UserPasswordResetTokenGet(rw web.ResponseWriter, req *web.Request) {
	// get email from url params

	queryMap := req.URL.Query()
	emailValues, emailOk := queryMap["email"]

	if !emailOk || len(emailValues) <= 0 {
		model := models.NewErrorResponse(constants.APIParsingQueryParams, models.NewAZError("Could not find email in url"), "Query params error")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	email := emailValues[0]
	fmt.Printf("test email: %s\n", email)
	var user models.User
	user.Email = email

	if err := c.DAL.GetUserByEmail(&user); err != nil {
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			// no such user
			c.Render(constants.StatusNotFound, nil, rw, req)
			return
		}
		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Error retrieving user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	fmt.Printf("test user reset token: %v\n", *user.ResetToken)

	if helpers.IsZeroString(user.ResetToken) {
		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError("reset token was empty"), "Error user has no token")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	model := models.TestResetToken{Token: *user.ResetToken}
	c.Render(constants.StatusOK, model, rw, req)
}

// UserPasswordResetTokenDelete deletes a users reset password token
//
//   DELETE /test/users/ResourcePasswordReset/userid
//
// Returns
//   204 No Content
func (c *TestContext) UserPasswordResetTokenDelete(rw web.ResponseWriter, req *web.Request) {
	params := req.PathParams
	userIDStr := params["user_id"]

	// TODO: determine if we can simply
	// pass in a partial struct with the things we
	// want to update, or do we need
	// to have custom DAL calls
	// that do this: _, err := db.Model(&book).Set("title = ?title").Returning("*").Update() (only updates title)
	// if we can pass in the struct, will mgo be ok with that?
	// the disadvantage of the non struct way is that the sql statements don't need to be kept up to date
	// as they will be written as per the struct
	// it will need a custom struct anyways, so for mgo it might not work...?
	// looks like mgo will handle partial structs just fine
	// we will need to create them however as separate structs for the DAL

	var user models.User

	user.ID = userIDStr

	// TODO: how would we set NIL to a time field?
	// ,sql should marshal this as NULL (?) => how would we set the 0 time? lol
	//user.ResetToken = nil
	//user.ResetTokenExpiry = nil

	// this will now handle three errors: id is not correct, user is not found, and some other error
	if err := c.DAL.ClearUserResetToken(&user); err != nil {
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			// no such user
			c.Render(constants.StatusNotFound, nil, rw, req)
			return
		}

		model := models.NewErrorResponse(constants.APIDatabaseUpdateUser, models.NewAZError(err.Error()), "Error updating user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	c.Render(constants.StatusNoContent, nil, rw, req)
}

// UserDelete deletes a user
//
//   DELETE /test/users/userid
//
// Returns
//   204 No Content
func (c *TestContext) UserDelete(rw web.ResponseWriter, req *web.Request) {
	// get user id
	params := req.PathParams
	userIDStr := params["user_id"]

	// just call delete
	var user models.User

	user.ID = userIDStr

	if err := c.DAL.DeleteUser(&user); err != nil {
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			// no such user
			c.Render(constants.StatusNotFound, nil, rw, req)
			return
		}
		model := models.NewErrorResponse(constants.APIDatabaseDeleteUser, models.NewAZError(err.Error()), "Error deleting user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}
	c.Render(constants.StatusNoContent, nil, rw, req)
}

// InvitationsDelete deletes all invitations
//
//   DELETE /test/users/invitations
//
// Returns
//   204 No Content
func (c *TestContext) InvitationsDelete(rw web.ResponseWriter, req *web.Request) {
	invites, err := c.DAL.GetAllInvitations()
	if err != nil {
		model := models.NewErrorResponse(constants.APIDatabaseDeleteUser, models.NewAZError(err.Error()), "Error deleting invites")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}
	for _, invite := range invites {
		if err := c.DAL.DeleteInvitation(invite); err != nil {
			model := models.NewErrorResponse(constants.APIDatabaseDeleteUser, models.NewAZError(err.Error()), "Error deleting invites")
			c.Render(constants.StatusInternalServerError, model, rw, req)
			return
		}
	}
	c.Render(constants.StatusNoContent, nil, rw, req)
}
