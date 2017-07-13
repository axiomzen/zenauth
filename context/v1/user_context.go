package v1

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/data"
	"github.com/axiomzen/zenauth/email"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/models"
	"github.com/gocraft/web"

	"github.com/twinj/uuid"
)

// UserContext for user authenticated routes
type UserContext struct {
	*APIAuthContext

	UserID string
}

var versionRegexp = regexp.MustCompile(`^/v[\d]+/`)

// AuthRequired Middleware: Authorizes a user by authenticating the Json Web Token
func (c *UserContext) AuthRequired(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes, Token: r.Header.Get(c.Config.AuthTokenHeader)}

	jwtTokenResult := jwt.Validate(c.Config.JwtClaimUserID)
	//fmt.Println("jwtTokenResult: " + jwtTokenResult.Message)

	switch jwtTokenResult.Status {
	case helpers.JWTokenStatusValid:

		if _, err := uuid.Parse(jwtTokenResult.Value); err != nil {
			c.UnauthorizedHandler(w, r)
			return
		}
		c.UserID = jwtTokenResult.Value

		c.Log = c.Log.WithField("userID", jwtTokenResult.Value)
		next(w, r)
	case helpers.JWTokenStatusExpired:
		c.ExpiredHandler(w, r)
	case helpers.JWTokenStatusInvalid, helpers.JWTokenNotAvailableYet:
		c.UnauthorizedHandler(w, r)
	}
}

// NewAuthToken creates a new auth token for a given admin id
func (c *UserContext) NewAuthToken(ID string) (*helpers.JWTHelper, error) {
	claims := make(map[string]interface{}, 2)
	claims[c.Config.JwtClaimUserID] = ID
	jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes}
	err := jwt.Generate(claims, c.Config.JwtUserTokenDuration)
	return &jwt, err
}

// renderUserResponseWithNewToken will render a UserResponse with a new token, given a user and a status
func (c *UserContext) renderUserResponseWithNewToken(user *models.User, status constants.HTTPStatusCode, sendVerificationEmail bool, w web.ResponseWriter, r *web.Request) {

	// create a new token

	tokenHelper, tokenErr := c.NewAuthToken(user.ID)

	if tokenErr != nil {
		model := models.NewErrorResponse(constants.APIAuthTokenCreation, models.NewAZError(tokenErr.Error()), "Could not create auth token")
		c.Render(constants.StatusInternalServerError, model, w, r)
		return
	}

	// TODO: add in JTI field (TODO: this could just be a stand alone table that we 'go' )
	// compound primary key
	// but again not sure what the use case is ( stolen/hacked account? )
	// if that was the case we should be logging anytime you hit a login/auth route
	// so we have more evidence

	// render response
	user.AuthToken = tokenHelper.Token

	// special case for 201 created
	// TODO: a more elegant way of doing this
	if status == constants.StatusCreated {
		// append the id of the user
		var b bytes.Buffer
		b.WriteString(r.URL.String())
		b.WriteString("/")

		b.WriteString(user.ID)

		w.Header().Set("Location", b.String())
	}

	c.Render(status, user, w, r)

	// TODO
	if sendVerificationEmail && user.Email != "" {
		// create the user email validation email
		// send the email
		// TODO: do we even have the route to verify the email for the user?
		// TODO: can/should this token have an expiry?
		// put email in it for claims
		emailer, err := email.Get(c.Config)
		if err != nil {
			model := models.NewErrorResponse(constants.APIVerifyEmailMessageError, models.NewAZError(err.Error()), "unable to generate verification email")
			c.Render(constants.StatusInternalServerError, model, w, r)
			return
		}

		claims := make(map[string]interface{}, 1)
		claims[c.Config.JwtClaimUserEmail] = user.Email
		jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes}
		if err := jwt.Generate(claims, c.Config.PasswordResetValidTokenDuration); err != nil {
			model := models.NewErrorResponse(constants.APIVerifyEmailMessageError, models.NewAZError(err.Error()), "unable to generate verification token")
			c.Render(constants.StatusInternalServerError, model, w, r)
			return
		}
		user.VerifyEmailToken = jwt.Token

		msg, err := email.GetVerifyEmailMessage(c.Config, user)
		if err != nil {
			model := models.NewErrorResponse(constants.APIVerifyEmailMessageError, models.NewAZError(err.Error()), "unable to generate verification email")
			c.Render(constants.StatusInternalServerError, model, w, r)
			return
		}

		// we don't care about email fails do we? perhaps log it
		go func(m *email.Message) {
			err := emailer.Send(m)
			if err != nil {
				c.Log.WithError(err).Warn("error sending email")
			}
		}(msg)
		// 	go models.CreateUserEmailValidationToken(user.Email)
	}
}

// VerifyEmail this link is clicked from a web browser (email client) and allows the user to verify
// their email (we sent the email with the link)
//
//   PUT /verify_email
// Returns
//   200 OK
func (c *UserContext) VerifyEmail(rw web.ResponseWriter, req *web.Request) {
	// so they will have a token in the url, and an email
	// tokens for this probably never need to expire or
	// be consumed, as all that it does is verify an email
	// so the check can be stateless
	queryMap := req.URL.Query()
	tokenSlice, tokenOk := queryMap["token"]
	emailSlice, emailOK := queryMap["email"]

	if !tokenOk || !emailOK {
		// TODO: this url is sent via a web browser
		// so I imagine they would want a better styled response
		// instead of an object
		msg := models.Message{Message: "400 - Bad Request (Missing params)"}
		c.Render(constants.StatusBadRequest, &msg, rw, req)
		return
	}

	// lower email
	emailAddr := strings.ToLower(emailSlice[0])

	jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes, Token: tokenSlice[0]}
	jwtTokenResult := jwt.Validate(c.Config.JwtClaimUserEmail)

	switch jwtTokenResult.Status {
	case helpers.JWTokenStatusValid:
		// verify that the emails are the same
		if jwtTokenResult.Value == emailAddr {
			// everything ok, set this users email to verified = true
			// we're not going to know to look by email, so we need a custom
			// method
			var user models.User
			user.Email = emailAddr
			user.Verified = true
			if err := c.DAL.UpdateUserVerified(&user); err != nil {
				// render error response
				msg := models.Message{Message: "500 - Bad Request (Database)"}
				c.Render(constants.StatusInternalServerError, &msg, rw, req)
				return
			}
			// render OK, or redirect?
			//************TODO: decide what you want to do here (redirect or what)
			c.Render(constants.StatusOK, &user, rw, req)
			return
		}
		// render error, email doesn't matches
		msg := models.Message{Message: "400 - Email doesn't match"}
		c.Render(constants.StatusBadRequest, &msg, rw, req)
	case helpers.JWTokenStatusExpired:
		// render expired
		msg := models.Message{Message: "400 - Token Exipired"}
		c.Render(constants.StatusBadRequest, &msg, rw, req)
	case helpers.JWTokenStatusInvalid, helpers.JWTokenNotAvailableYet:
		msg := models.Message{Message: "400 - Invalid Token"}
		c.Render(constants.StatusBadRequest, &msg, rw, req)
	}
}

// ForgotPassword route
//
//   PUT /forgot_password?email=:email:
//
// Expects one query param, email
// This is step 1 of 3
//
// This was the decided design from timeline, but there are alternatives,
// including simply changing the password on the users behalf (to something random)
// and sending the new password in the email
//
// Returns
//   204 No Content
func (c *UserContext) ForgotPassword(rw web.ResponseWriter, req *web.Request) {

	queryMap := req.URL.Query()
	emailSlice, emailOk := queryMap["email"]

	if !emailOk {
		model := models.NewErrorResponse(constants.APIParsingQueryParams, models.NewAZError("email expected"), "query parameter missing")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	// generate a JWT with expiry and other claims (email) so that we don't have to check the DB on the get
	claims := make(map[string]interface{}, 2)
	// TODO: why did we do this?
	//email := strings.Replace(emailSlice[0], " ", "+", -1)
	// TODO: test email with spaces
	claims[c.Config.JwtClaimUserEmail] = emailSlice[0]
	jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes}
	err := jwt.Generate(claims, c.Config.PasswordResetValidTokenDuration)

	if err != nil {
		model := models.NewErrorResponse(constants.APIParsingQueryParams, models.NewAZError(err.Error()), "unable to generate reset token")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	// save the token in the database, so that it is single use only
	// might as well save it to the users table, as when we consume it
	// we need to update the users password, so we need to do a where condition
	// anyways

	var user models.User
	//fmt.Printf("user email: %s\n", emailSlice[0])
	user.Email = emailSlice[0]
	user.ResetToken = &jwt.Token
	//fmt.Printf("reset token before: %s\n", user.ResetToken.String)

	if err = c.DAL.CreateUserResetToken(&user); err != nil {
		// check to see what kind of error; if its none affected just return no content anyways
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			// this email was invalid, but we don't allow them to know that
			// and we don't send emails to random email addresses
			c.Render(constants.StatusNoContent, nil, rw, req)
			return
		}
		// some other error
		model := models.NewErrorResponse(constants.APIParsingQueryParams, models.NewAZError(err.Error()), "unable to save reset token")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	//fmt.Printf("reset token after: %s\n", user.ResetToken.String)

	// send the reset password email with the generated token
	emailer, err := email.Get(c.Config)
	if err != nil {
		model := models.NewErrorResponse(constants.APIForgotPasswordMessageError, models.NewAZError(err.Error()), "unable to generate reset token email")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}
	msg, err := email.GetResetPasswordMessage(c.Config, &user)
	if err != nil {
		model := models.NewErrorResponse(constants.APIForgotPasswordMessageError, models.NewAZError(err.Error()), "unable to generate reset token email")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	go func(m *email.Message) {
		err := emailer.Send(m)
		if err != nil {
			c.Log.WithError(err).Warn("error sending email")
		}
	}(msg)
	// render response
	c.Render(constants.StatusNoContent, nil, rw, req)
}

// ResetPassword route
//
//   POST /reset_password
//
// This is step 3 of 3
//
// Returns
//   200 Status OK
func (c *UserContext) ResetPassword(rw web.ResponseWriter, req *web.Request) {
	// so at this point they have been sent to the page
	// and we are still using that existing token
	// which hasn't been consumed yet
	// the code will automatically check for time expired
	// so the above route seems really unnessesary
	// as by the time they pick a new password it could be expired already

	// after all this, we still need to check the new password to see if its legit
	// (might be wise to do PW checks on the front end to avoid hitting our service)
	// TODO: decide where the token will go (request body? header? url?)

	var userPasswordReset models.UserPasswordReset
	// decode request
	if !c.DecodeHelper(&userPasswordReset, "Couldn't decode userPasswordReset", rw, req) {
		return
	}

	// verify the token
	jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes, Token: userPasswordReset.Token}
	jwtTokenResult := jwt.Validate(c.Config.JwtClaimUserEmail)

	switch jwtTokenResult.Status {
	case helpers.JWTokenStatusValid:
		// verify that the emails are the same
		if jwtTokenResult.Value != userPasswordReset.Email {
			// render error, email doesn't matches
			msg := models.Message{Message: "400 - Email doesn't match"}
			c.Render(constants.StatusBadRequest, &msg, rw, req)
			return
		}

		if len(userPasswordReset.NewPassword) < int(c.Config.MinPasswordLength) {
			model := models.NewErrorResponse(constants.APIValidationPasswordTooShort,
				models.NewAZError(fmt.Sprintf("Password needs to be at least %d characters long!", c.Config.MinPasswordLength)), "Could not create account")
			c.Render(constants.StatusBadRequest, model, rw, req)
			return
		}

		var user models.User
		user.ResetToken = &userPasswordReset.Token
		user.Email = userPasswordReset.Email

		newHash, hashErr := helpers.HashPasswordBcrypt(userPasswordReset.NewPassword, int(c.Config.BcryptCost))

		if hashErr != nil {
			model := models.NewErrorResponse(constants.APIParsingPasswordHash, models.NewAZError(hashErr.Error()), "Could not update user")
			c.Render(constants.StatusInternalServerError, model, rw, req)
			return
		}

		// set the hash
		user.Hash = &newHash
		//fmt.Printf("new user hash: %s\n", *user.Hash)

		// this will set the token to null, and update the password hash for the user by email
		// only if the token matches
		if err := c.DAL.ConsumeUserResetToken(&user); err != nil {
			dalErr, _ := err.(data.DALError)
			// no such user/email
			if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
				c.NotFound(rw, req)
				return
			}

			model := models.NewErrorResponse(constants.APIParsingPasswordHash, models.NewAZError(err.Error()), "Could not update user")
			c.Render(constants.StatusInternalServerError, model, rw, req)
			return
		}

		c.renderUserResponseWithNewToken(&user, constants.StatusOK, false, rw, req)

		// TODO: localization (we need to get a string via id => it will have appropriate %s etc)
		// TODO: does one instance have 1 file or do all instances have all files and a localization?

	case helpers.JWTokenStatusExpired:
		// render expired
		msg := models.Message{Message: "400 - Reset Request Expired"}
		c.Render(constants.StatusBadRequest, &msg, rw, req)
	case helpers.JWTokenStatusInvalid, helpers.JWTokenNotAvailableYet:
		msg := models.Message{Message: "400 - Invalid Token"}
		c.Render(constants.StatusBadRequest, &msg, rw, req)
	}
}

// Exists route - expects one query param, email
//
//   GET /exists?email=:email:
//
// Returns
//   200 OK
func (c *UserContext) Exists(rw web.ResponseWriter, req *web.Request) {
	// check for email existance
	queryMap := req.URL.Query()
	response := models.Exists{}
	_, ok := queryMap["email"]

	if !ok {
		model := models.NewErrorResponse(constants.APIParsingQueryParams, models.NewAZError("email expected"), "query parameter missing")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}
	// fixup email
	email := strings.Replace(queryMap.Get("email"), " ", "+", -1)
	email = strings.ToLower(strings.Trim(email, " "))

	// check to see if user exists
	var user models.User
	user.Email = email
	response.Exists = c.DAL.GetUserByEmail(&user) == nil

	// render resposne
	c.Render(constants.StatusOK, &response, rw, req)
}

// Get route - Returns the current user information
//
//   GET /users
//
// Returns
//   200 OK
func (c *UserContext) GetSelf(rw web.ResponseWriter, req *web.Request) {

	var user models.User
	user.ID = c.UserID

	// get user
	if err := c.DAL.GetUserByID(&user); err != nil {
		dalErr, _ := err.(data.DALError)
		// no such user/email
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			c.NotFound(rw, req)
			return
		}

		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Could not get user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	// get token from header
	user.AuthToken = req.Header.Get(c.Config.AuthTokenHeader)

	// render response
	c.Render(constants.StatusOK, user, rw, req)
}

// Get route - Returns a user information if the user is registered OR invited
//
//   GET /users/:ID
//
// Returns
//   200 OK
func (c *UserContext) Get(rw web.ResponseWriter, req *web.Request) {

	var user models.User
	user.ID = req.PathParams["id"]

	// get user
	if err := c.DAL.GetUserByID(&user); err != nil {
		dalErr, _ := err.(data.DALError)
		// no such user/email
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			if c.renderInvitation(rw, req, user.ID) {
				// If the invitation got rendered, we're done.
				return
			}
			c.NotFound(rw, req)
			return
		}

		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Could not get user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}
	view, err := user.ProtobufPublic()
	if err != nil {
		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Could not generate view for the user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	// render response
	c.Render(constants.StatusOK, view, rw, req)
}

func (c *UserContext) renderInvitation(rw web.ResponseWriter, req *web.Request, invitationID string) bool {
	invitation := models.Invitation{
		ID: invitationID,
	}
	if err := c.DAL.GetInvitationByID(&invitation); err != nil {
		return false
	}
	view, err := invitation.UserPublicProtobuf()
	if err != nil {
		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Could not generate view for the user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return true
	}
	c.Render(constants.StatusOK, view, rw, req)
	return true
}

// PasswordPut changes the user password
//
//   PUT /password
//
// In case someone gets on their computer we need the old password
// so we need a special route
//
// Returns
//   200 OK
func (c *UserContext) PasswordPut(rw web.ResponseWriter, req *web.Request) {
	var userChangePassword models.UserChangePassword
	// decode request
	if !c.DecodeHelper(&userChangePassword, "Couldn't decode userChangePassword", rw, req) {
		return
	}

	// TODO: refactor
	// check new password
	if len(userChangePassword.NewPassword) < int(c.Config.MinPasswordLength) {
		model := models.NewErrorResponse(constants.APIValidationPasswordTooShort,
			models.NewAZError(fmt.Sprintf("Password needs to be at least %d characters long!", c.Config.MinPasswordLength)), "Could not update password")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	// TODO: check old password
	// we can't just have a where clause as
	// we don't store the old password in the DB

	// get the user
	var user models.User
	user.ID = c.UserID
	if err := c.DAL.GetUserByID(&user); err != nil {
		// no such user
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			c.NotFound(rw, req)
			return
		}
		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Could not get user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	// if we have no hash, then we have no old password
	if helpers.IsZeroString(user.Hash) {
		email := "NULL"
		if user.Email != "" {
			email = user.Email
		}
		model := models.NewErrorResponse(constants.APIDatabaseUpdateUser,
			models.NewAZError("No password associated with this email: "+email), "Could not update user")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	// check to see if old password is correct
	//fmt.Println("user hash before checkpassword: %s\n", *user.Hash)
	//fmt.Println("user password before checkpassword: %s\n", userChangePassword.OldPassword)

	if passwordOK, err := helpers.CheckPasswordBcrypt(*user.Hash, userChangePassword.OldPassword); err != nil {

		model := models.NewErrorResponse(constants.APIDatabaseUpdateUser, models.NewAZError(err.Error()), "Could not check user password")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	} else if !passwordOK {
		// wrong password
		model := models.NewErrorResponse(constants.APIDatabaseUpdateUser, models.NewAZError("Old password incorrect"), "Could not update user")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	// safe to update has now
	// generate new hash

	newHash, hashErr := helpers.HashPasswordBcrypt(userChangePassword.NewPassword, int(c.Config.BcryptCost))

	if hashErr != nil {
		model := models.NewErrorResponse(constants.APIParsingPasswordHash, models.NewAZError(hashErr.Error()), "Could not generate user hash")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	// attempt to update the password, may fail
	if err := c.DAL.UpdateUserHash(newHash, &user); err != nil {
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			c.NotFound(rw, req)
			return
		}
		model := models.NewErrorResponse(constants.APIDatabaseUpdateUser, models.NewAZError(err.Error()), "Could not update user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	// everything ok
	// get token from header
	// TODO: we could invalidate all the old tokens here
	// and generate a new one
	// TODO: generate new token and invalidate old one?
	user.AuthToken = req.Header.Get(c.Config.AuthTokenHeader)
	// render response
	c.Render(constants.StatusOK, user, rw, req)
}

// EmailPut changes the email address of the user
//
//   PUT /email
//
// TODO: reevaluate why do we need a special endpoint
// as apposed to just build this into the user PUT?
//
// Returns
//   200 OK
func (c *UserContext) EmailPut(rw web.ResponseWriter, req *web.Request) {
	var userChangeEmail models.UserChangeEmail

	// decode request
	if !c.DecodeHelper(&userChangeEmail, "Couldn't decode UserChangeEmail", rw, req) {
		return
	}

	// fix up
	userChangeEmail.Email = strings.ToLower(strings.Trim(userChangeEmail.Email, " "))

	// validate new email
	if strings.Count(userChangeEmail.Email, "@") != 1 {
		model := models.NewErrorResponse(constants.APIValidationEmailNotValid, models.NewAZError("Please enter a valid email address"), "Could not create account")
		c.Render(constants.StatusBadRequest, model, rw, req)
		return
	}

	userChangeEmail.ID = c.UserID
	var user models.User
	//user.ID = c.UserID
	//user.Email = &userChangeEmail.Email

	// TODO: could also just use UpdateUser
	//if err := c.DAL.UpdateUserEmail(&user); err != nil {
	if err := c.DAL.UpdateUser(&userChangeEmail, &user); err != nil {
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			c.NotFound(rw, req)
			return
		} else if dalErr.ErrorCode == data.DALErrorCodeUniqueEmail {
			model := models.NewErrorResponse(constants.APIDatabaseUpdateUser, models.NewAZError(err.Error()), "Could not update user")
			c.Render(constants.StatusBadRequest, model, rw, req)
			return
		}
		model := models.NewErrorResponse(constants.APIDatabaseUpdateUser, models.NewAZError(err.Error()), "Could not update user")
		c.Render(constants.StatusInternalServerError, model, rw, req)
		return
	}

	// everything ok
	// get token from header
	// TODO: we could invalidate all the old tokens here
	// and generate a new one
	user.AuthToken = req.Header.Get(c.Config.AuthTokenHeader)
	// render response
	c.Render(constants.StatusOK, user, rw, req)
}

// Login logs a user in
//
//   POST /login
//
// Returns
//   200 OK
func (c *UserContext) Login(w web.ResponseWriter, req *web.Request) {

	token := req.Header.Get(c.Config.AuthTokenHeader)
	if token != "" {
		jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes, Token: token}
		jwtTokenResult := jwt.Validate(c.Config.JwtClaimUserID)

		if jwtTokenResult.Status == helpers.JWTokenStatusValid {

			if _, err := uuid.Parse(jwtTokenResult.Value); err == nil {
				c.UserID = jwtTokenResult.Value

				c.Log = c.Log.WithField("userID", jwtTokenResult.Value)

				// get the user
				var user models.User
				user.ID = c.UserID
				if err := c.DAL.GetUserByID(&user); err != nil {
					dalErr, _ := err.(data.DALError)
					if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
						model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "User does not exist")
						c.Render(constants.StatusUnauthorized, model, w, req)
						return
					}
					model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Could not get user")
					c.Render(constants.StatusInternalServerError, model, w, req)
					return
				}
				// just send them back their user with same token
				user.AuthToken = token
				c.Render(constants.StatusOK, user, w, req)
				return
			}
		}
	}

	// token either invalid or missing
	var login models.Login

	// decode request
	if !c.DecodeHelper(&login, "Couldn't decode login", w, req) {
		return
	}

	// If no parsing error, prepare the response
	login.Email = strings.ToLower(strings.Trim(login.Email, " "))

	var user models.User
	user.Email = login.Email
	//fmt.Println("user email: " +  user.Email.String)
	user.UserName = login.UserName
	var err error
	if c.Config.RequireUsername {
		err = c.DAL.GetUserByEmailOrUserName(&user)
	} else {
		err = c.DAL.GetUserByEmail(&user)
	}
	if err != nil {
		dalErr, _ := err.(data.DALError)
		if dalErr.ErrorCode == data.DALErrorCodeNoneAffected {
			model := models.NewErrorResponse(constants.APILoginSignupInvalidCombination, models.NewAZError(err.Error()), "Invalid email/username/password combination")
			c.Render(constants.StatusUnauthorized, model, w, req)
			return
		}
		model := models.NewErrorResponse(constants.APIDatabaseGetUser, models.NewAZError(err.Error()), "Could not get user")
		c.Render(constants.StatusInternalServerError, model, w, req)
		return
	}

	// *********TODO: Optionally check for verified emails here
	// ********* we should wrap this in an .EmailVerification variable

	// if !user.Verified {
	// 	model := models.NewErrorResponse(constants.APILoginNotVerified, models.NewAZError("email address not verified"), "User must validate their email first")
	// 	c.Render(constants.StatusUnauthorized, model, w, req)
	// 	return
	// }

	// check that they have a password - not sure how they wouldn't
	if helpers.IsZeroString(user.Hash) {
		model := models.NewErrorResponse(constants.APILoginSignupInvalidCombination, models.NewAZError("No password associated with this email: "+user.Email), "Invalid email/username/password combination")
		c.Render(constants.StatusBadRequest, model, w, req)
		return
	}

	//fmt.Printf("User hash (login): %s\n", *user.Hash)
	//fmt.Printf("User password (login): %s\n", login.Password)

	if passwordOK, err := helpers.CheckPasswordBcrypt(*user.Hash, login.Password); err != nil {

		model := models.NewErrorResponse(constants.APIDatabaseUpdateUser, models.NewAZError(err.Error()), "Could not check user password")
		c.Render(constants.StatusInternalServerError, model, w, req)
		return
	} else if !passwordOK {
		// wrong password
		model := models.NewErrorResponse(constants.APILoginSignupInvalidCombination, models.NewAZError("Username/Password combination incorrect"), "Invalid email/username/password combination")
		c.Render(constants.StatusUnauthorized, model, w, req)
		return
	}

	go func(user models.User, login models.Login, c *UserContext) {
		// Check to see if user hash meets current security standards
		if update, newHash := helpers.UpgradeHashBcrypt(*user.Hash, login.Password, c.Config.BcryptCost, c.Config.AllowHashDowngrades); update {

			// attempt to update the password, may fail
			if err := c.DAL.UpdateUserHash(newHash, &user); err != nil {
				c.Log.WithError(err).WithField("code", constants.APIDatabaseUpdate).Error("Could not update user hash")
			}

		}
	}(user, login, c)

	c.renderUserResponseWithNewToken(&user, constants.StatusOK, false, w, req)
}

// Signup signs a user up via email
//
//   POST /signup
//
// Assumes format:
//   {
//     "email":"example@email.ca",
//     "password":"aPassword1"
//   }
//
// Returns
//   201 Created
func (c *UserContext) Signup(w web.ResponseWriter, req *web.Request) {

	var signup models.Signup

	// decode request
	if !c.DecodeHelper(&signup, "Couldn't decode signup", w, req) {
		return
	}

	//check password long enough
	if len(signup.Password) < int(c.Config.MinPasswordLength) {
		model := models.NewErrorResponse(constants.APIValidationPasswordTooShort,
			models.NewAZError(fmt.Sprintf("Password needs to be at least %d characters long!", c.Config.MinPasswordLength)), "Could not create account")
		c.Render(constants.StatusBadRequest, model, w, req)
		return
	}

	// check email
	// TODO: perhaps refactor
	if strings.Count(signup.Email, "@") == 0 {
		model := models.NewErrorResponse(constants.APIValidationEmailNotValid, models.NewAZError("Please enter a valid email address"), "Could not create account")
		c.Render(constants.StatusBadRequest, model, w, req)
		return
	}

	if c.Config.RequireUsername && signup.UserName == "" {
		model := models.NewErrorResponse(constants.APIValidationUserNameNotValid, models.NewAZError("Please enter a username"), "Could not create account")
		c.Render(constants.StatusBadRequest, model, w, req)
		return
	}

	// lower case the email
	signup.Email = helpers.EmailSanitize(signup.Email)

	// create a user var
	user := models.User{}
	user.Email = signup.Email
	user.UserName = signup.UserName

	// Verify that no user with this email exists
	if c.Config.RequireUsername {
		if err := c.DAL.GetUserByEmailOrUserName(&user); err == nil {
			model := models.NewErrorResponse(constants.APIDatabaseGetUser,
				models.NewAZError("User with email/username already exists"), "Email or Username already in use/exists")
			c.Render(constants.StatusForbidden, model, w, req)
			return
		}
	} else {
		if err := c.DAL.GetUserByEmail(&user); err == nil {
			model := models.NewErrorResponse(constants.APIDatabaseGetUser,
				models.NewAZError("User with email already exists"), "Email already in use/exists")
			c.Render(constants.StatusForbidden, model, w, req)
			return
		}
	}

	hash, hashErr := helpers.HashPasswordBcrypt(signup.Password, int(c.Config.BcryptCost))

	if hashErr != nil {
		model := models.NewErrorResponse(constants.APIParsingPasswordHash, models.NewAZError(hashErr.Error()), "Could not create new user")
		c.Render(constants.StatusBadRequest, model, w, req)
		return
	}

	user.Hash = &hash

	userErr := c.DAL.CreateUser(&user)

	if userErr != nil {

		dalErr, _ := userErr.(data.DALError)

		if dalErr.ErrorCode == data.DALErrorCodeUniqueEmail {
			model := models.NewErrorResponse(constants.APIEmailInUse, models.NewAZError(userErr.Error()), "Email already in use/exists")
			c.Render(constants.StatusForbidden, model, w, req)
			return
		}

		model := models.NewErrorResponse(constants.APIDatabaseCreateUser, models.NewAZError(userErr.Error()), "Could not create new user")
		c.Render(constants.StatusBadRequest, model, w, req)
		return
	}

	//********* TODO: if you want the email to be verified first, you would need
	//********* to not give them an auth token here
	//********* the current logic allows signup without email verification
	//********* TODO have a .EmailVerification option for hatch which nicely wraps
	//********* this functionality

	// render a user response
	c.renderUserResponseWithNewToken(&user, constants.StatusCreated, true, w, req)
}
