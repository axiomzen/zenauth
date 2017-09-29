package v1

import (
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/data"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/models"
	"github.com/gocraft/web"
	//"fmt"
)

// FacebookContext for Facebook login support
type FacebookContext struct {
	*UserContext
}

// validateFacebookUser helper function
func (c *FacebookContext) validateFacebookUser(fbUser *models.FacebookUser, rw web.ResponseWriter, req *web.Request) bool {
	// TODO Change to check hashed password against db & require username and password fields
	if fbUser.FacebookID == "" || fbUser.FacebookToken == "" {
		model := models.NewErrorResponse(constants.APIValidation, models.NewAZError("Missing a field in request"), "Error with request")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return false
	}

	if valid, err := helpers.ValidateFacebookLogin(fbUser.FacebookID, fbUser.FacebookToken, c.Config.FacebookAppID, c.Config.FacebookAppSecret); err != nil {
		model := models.NewErrorResponse(constants.APIFacebookLoginNotValid, models.NewAZError(err.Error()), "Error with fb login request")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return false
	} else if !valid {
		model := models.NewErrorResponse(constants.APIFacebookLoginNotValid, models.NewAZError(err.Error()), "Could not validate facebook token")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return false
	}
	fbAPIUser, err := helpers.GetFacebookUserInfo(
		fbUser.FacebookID,
		fbUser.FacebookToken,
		c.Config.FacebookAppID,
		c.Config.FacebookAppSecret)
	if err != nil {
		c.Log.WithError(err).Error("Could not retreive user profile picture")
	} else {
		fbUser.FacebookPicture = fbAPIUser.ProfilePicture
	}
	return true
}

// createFacebookUser helper function
func (c *FacebookContext) createFacebookUser(user *models.User, rw web.ResponseWriter, req *web.Request) bool {
	if err := c.DAL.CreateUser(user); err != nil {
		// facebook id might not be unique
		// email might not be unique
		dalErr, _ := err.(data.DALError)

		switch dalErr.ErrorCode {
		case data.DALErrorCodeUniqueEmail:
			model := models.NewErrorResponse(constants.APIEmailInUse, models.NewAZError(err.Error()), "Email already in use/exists")
			c.Render(constants.StatusForbidden, model, rw, req)
			return false
		case data.DALErrorCodeFacebookIDUnique:
			model := models.NewErrorResponse(constants.APISocialAccountExists, models.NewAZError(err.Error()), "Social account already exists")
			c.Render(constants.StatusForbidden, model, rw, req)
			return false
		default:
			model := models.NewErrorResponse(constants.APIDatabaseCreateUser, models.NewAZError(err.Error()), "Could not create new User")
			c.Render(constants.StatusInternalServerError, model, rw, req)
			return false
		}
	}
	return true
}

// Login logs a user in (via facebook)
//   POST /fblogin
//
// Returns
//   200 OK
func (c *FacebookContext) Login(rw web.ResponseWriter, req *web.Request) {

	var fbLogin models.FacebookUser

	// decode request
	if !c.DecodeHelper(&fbLogin, "Couldn't decode facebook user", rw, req) {
		return
	}

	if !c.validateFacebookUser(&fbLogin, rw, req) {
		return
	}

	var user models.User

	user.FacebookUser = fbLogin

	if err := c.DAL.GetUserByFacebookID(&user); err != nil {
		model := models.NewErrorResponse(constants.APILoginUserDoesNotExist, models.NewAZError(err.Error()), "User does not exist")
		c.Render(constants.StatusForbidden, model, rw, req)
		return
	}

	c.renderUserResponseWithNewToken(&user, constants.StatusOK, false, rw, req)
}

// Signup signs a user up via facebook (optionally links their account to an existing account)
//
//   POST /fbsignup
//
// Returns
//   201 Created
func (c *FacebookContext) Signup(rw web.ResponseWriter, req *web.Request) {

	var fbSignup models.FacebookSignup

	// decode request
	if !c.DecodeHelper(&fbSignup, "Couldn't decode facebook signup", rw, req) {
		return
	}

	// validate
	if !c.validateFacebookUser(&fbSignup.FacebookUser, rw, req) {
		return
	}

	// create a user
	user := models.User{}
	user.FacebookUser = fbSignup.FacebookUser
	// user.Email = helpers.EmailSanitize(fbSignup.FacebookEmail)

	// otherwise, we want to create (should populate user with new id)
	if !c.createFacebookUser(&user, rw, req) {
		return
	}

	// render response
	c.renderUserResponseWithNewToken(&user, constants.StatusCreated, true, rw, req)
}

// Link links a facebook account to an existing account
//
//   POST /users/fblink
//
// Returns
//   200 OK
func (c *FacebookContext) Link(rw web.ResponseWriter, req *web.Request) {

	var fbUpdate models.FacebookUpdate

	// decode request
	if !c.DecodeHelper(&fbUpdate, "Couldn't decode facebook update", rw, req) {
		return
	}

	// validate
	if !c.validateFacebookUser(&fbUpdate.FacebookUser, rw, req) {
		return
	}

	// user id will be populated at this point from the token
	fbUpdate.ID = c.UserID
	var user models.User

	if err := c.DAL.UpdateUser(&fbUpdate, &user); err != nil {
		dalErr, _ := err.(data.DALError)

		if dalErr.ErrorCode == data.DALErrorCodeFacebookIDUnique {
			model := models.NewErrorResponse(constants.APISocialAccountExists, models.NewAZError(err.Error()), "Social account already exists")
			c.Render(constants.StatusForbidden, model, rw, req)
			return
		} else if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			c.NotFound(rw, req)
			return
		}

		// user might not be found
		model := models.NewErrorResponse(constants.APIDatabaseCreateUser, models.NewAZError(err.Error()), "Could not update User")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}
	// create a new token
	c.renderUserResponseWithNewToken(&user, constants.StatusOK, false, rw, req)
}

// Facebook will log the user in, and create an account if it doesn't already exist
//
//   POST /facebook
//
// Returns
//   201 Created
func (c *FacebookContext) Facebook(rw web.ResponseWriter, req *web.Request) {

	var fbSignup models.FacebookSignup

	// decode request
	if !c.DecodeHelper(&fbSignup, "Couldn't decode facebook signup", rw, req) {
		return
	}

	// validate
	if !c.validateFacebookUser(&fbSignup.FacebookUser, rw, req) {
		return
	}

	// create a user
	user := models.User{}
	user.FacebookUser = fbSignup.FacebookUser

	if err := c.DAL.UpdateUserFacebookInfo(&user); err == nil {
		c.renderUserResponseWithNewToken(&user, constants.StatusOK, false, rw, req)
		return
	}

	// Else signup
	// user.Email = helpers.EmailSanitize(fbSignup.FacebookEmail)

	// otherwise, we want to create (should populate user with new id)
	if !c.createFacebookUser(&user, rw, req) {
		return
	}

	// render response
	c.renderUserResponseWithNewToken(&user, constants.StatusCreated, true, rw, req)
}
