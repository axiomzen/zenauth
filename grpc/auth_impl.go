package grpc

import (
	"fmt"
	"strings"

	context "golang.org/x/net/context"

	"google.golang.org/grpc/metadata"

	"github.com/Sirupsen/logrus"
	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/data"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/models"
	"github.com/axiomzen/zenauth/protobuf"
	pEmpty "github.com/golang/protobuf/ptypes/empty"
	"github.com/twinj/uuid"
)

type Auth struct {
	Config *config.ZENAUTHConfig
	DAL    data.ZENAUTHProvider
	Log    *logrus.Entry
}

// GetCurrentUser implements the action to return the user from the session token.
func (auth *Auth) GetCurrentUser(ctx context.Context, _ *pEmpty.Empty) (*protobuf.User, error) {
	userID, err := auth.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	user.ID = userID

	// get user
	if err := auth.DAL.GetUserByID(&user); err != nil {
		return nil, err
	}

	return user.Protobuf()
}

// GetUserByID implements the action to return the user from the ID.
func (auth *Auth) GetUserByID(ctx context.Context, userID *protobuf.UserID) (*protobuf.UserPublic, error) {

	// Get the current user to make sure it's an authenticated request
	_, err := auth.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user models.User
	user.ID = userID.GetId()

	// get user
	if err := auth.DAL.GetUserByID(&user); err != nil {
		return nil, err
	}

	return user.ProtobufPublic()
}

// LinkUser implements the action to link a user
func (auth *Auth) LinkUser(ctx context.Context, invite *protobuf.InvitationCode) (*protobuf.UserPublic, error) {
	if !constants.InvitationTypes[invite.GetType()] {
		return nil, fmt.Errorf("Invitation type %s not supported", invite.GetType())
	}
	userID, err := auth.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	user.ID = userID
	// get user
	if err := auth.DAL.GetUserByID(&user); err != nil {
		return nil, err
	}

	// Check if we are just linking an invite
	invitation := models.Invitation{
		Code: invite.GetInviteCode(),
		Type: invite.GetType(),
	}
	if err := auth.DAL.GetInvitation(&invitation); err == nil {
		if userInfoUpdateErr := invitation.UpdateUserWithInvitationInfo(&user); userInfoUpdateErr != nil {
			return nil, userInfoUpdateErr
		} else if userInfoUpdateErr = auth.DAL.UpdateUser(&user, &user); userInfoUpdateErr != nil {
			return nil, userInfoUpdateErr
		} else if delInvErr := auth.DAL.DeleteInvitation(&invitation); delInvErr != nil {
			return nil, delInvErr
		}
		userPub, err := invitation.UserPublicProtobuf()
		userPub.Status = protobuf.UserStatus_merged
		return userPub, err
	} else if dalErr, ok := err.(data.DALError); !ok || dalErr.ErrorCode != data.DALErrorCodeNoneAffected {
		// Any error other than could not find the invitation
		return nil, err
	}
	// This means that no invitation exists, check for a user with this inv code

	var linkToUser models.User
	var linkUserErr error
	switch invite.GetType() {
	case constants.InvitationTypeEmail:
		user.Email = invite.GetInviteCode()
		linkToUser.Email = invite.GetInviteCode()
		linkUserErr = auth.DAL.GetUserByEmail(&linkToUser)
	case constants.InvitationTypeFacebook:
		user.FacebookID = invite.GetInviteCode()
		linkToUser.FacebookID = invite.GetInviteCode()
		linkUserErr = auth.DAL.GetUserByFacebookID(&linkToUser)
	default:
		// Should never get here as we check the invite type above,
		linkUserErr = fmt.Errorf("Invitation type %s not supported", invite.GetType())
	}

	if linkUserErr == nil {
		// User found, delete and return
		if mergeUserErr := auth.DAL.MergeUsers(&user, &linkToUser); mergeUserErr != nil {
			return nil, mergeUserErr
		}
		mergedUser, returnErr := linkToUser.ProtobufPublic()
		mergedUser.Status = protobuf.UserStatus_merged
		return mergedUser, returnErr
	}
	auth.Log.WithError(linkUserErr).Debug("Could not retrieve social account to link")

	if err := auth.DAL.UpdateUser(&user, &user); err != nil {
		return nil, err
	}
	return user.ProtobufPublic()
}

// GetUserByID implements the action to return the user from the ID.
func (auth *Auth) GetUsersByIDs(ctx context.Context, userIDs *protobuf.UserIDs) (*protobuf.UsersPublic, error) {

	// Get the current user to make sure it's an authenticated request
	_, err := auth.getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var users models.Users
	for _, id := range userIDs.GetIds() {
		users = append(users, &models.User{
			UserBase: models.UserBase{ID: id},
		})
	}

	// get users
	if err := auth.DAL.GetUsersByIDs(&users); err != nil {
		return nil, err
	}

	return users.ProtobufPublic()
}

// AuthUserByEmail implements the action to either signup or login
func (auth *Auth) AuthUserByEmail(ctx context.Context, emailAuth *protobuf.UserEmailAuth) (*protobuf.User, error) {

	user := models.User{
		UserBase: models.UserBase{
			Email:    helpers.EmailSanitize(emailAuth.GetEmail()),
			UserName: emailAuth.GetUserName(),
		},
	}
	var err error
	if auth.Config.RequireUsername {
		err = auth.DAL.GetUserByEmailOrUserName(&user)
	} else {
		err = auth.DAL.GetUserByEmail(&user)
	}
	if err == nil {
		// Can just login
		// check that they have a password - not sure how they wouldn't
		if helpers.IsZeroString(user.Hash) {
			return nil, fmt.Errorf("Wrong account type (No password saved)")
		}

		if passwordOK, err := helpers.CheckPasswordBcrypt(*user.Hash, emailAuth.GetPassword()); err != nil {
			return nil, err
		} else if !passwordOK {
			// wrong password
			return nil, fmt.Errorf("Invalid email/username/password combination")
		}
		authToken, tokenErr := auth.NewAuthToken(user.ID)
		if tokenErr != nil {
			return nil, tokenErr
		}
		user.AuthToken = authToken
		return user.Protobuf()
	}

	if dalErr, isDALError := err.(data.DALError); isDALError && dalErr.ErrorCode != data.DALErrorCodeNoneAffected {
		// Error getting user, not that it doesn't exist
		return nil, err
	}
	// Else, Sign Up
	// Validate
	if len(emailAuth.GetPassword()) < int(auth.Config.MinPasswordLength) {
		// check password long enough
		return nil, fmt.Errorf("Password too short")
	} else if strings.Count(user.Email, "@") == 0 {
		// check email
		return nil, fmt.Errorf("Invalid Email")
	} else if auth.Config.RequireUsername && user.UserName == "" {
		return nil, fmt.Errorf("Please enter a username")
	}

	hash, hashErr := helpers.HashPasswordBcrypt(emailAuth.GetPassword(), int(auth.Config.BcryptCost))

	if hashErr != nil {
		return nil, hashErr
	}

	user.Hash = &hash

	if userErr := auth.DAL.CreateUser(&user); userErr != nil {
		return nil, userErr
	}
	// Generate the auth token
	authToken, tokenErr := auth.NewAuthToken(user.ID)
	if tokenErr != nil {
		return nil, tokenErr
	}
	user.AuthToken = authToken
	return user.Protobuf()
}

// AuthUserByFacebook implements the action to return the user from the ID.
func (auth *Auth) AuthUserByFacebook(ctx context.Context, facebookAuth *protobuf.UserFacebookAuth) (*protobuf.User, error) {

	// validate
	if facebookAuth.GetFacebookID() == "" || facebookAuth.GetFacebookToken() == "" {
		return nil, fmt.Errorf("Missing a field in request")
	}

	if valid, err := helpers.ValidateFacebookLogin(facebookAuth.GetFacebookID(), facebookAuth.GetFacebookToken(), auth.Config.FacebookAppID, auth.Config.FacebookAppSecret); err != nil {
		return nil, err
	} else if !valid {
		return nil, fmt.Errorf("Could not validate facebook token")
	}
	// create a user
	user := models.User{
		FacebookUser: models.FacebookUser{
			FacebookID:       facebookAuth.GetFacebookID(),
			FacebookEmail:    facebookAuth.GetFacebookEmail(),
			FacebookUsername: facebookAuth.GetFacebookUsername(),
			FacebookToken:    facebookAuth.GetFacebookToken(),
		},
	}
	if err := auth.DAL.UpdateUserFacebookToken(&user); err == nil {
		authToken, tokenErr := auth.NewAuthToken(user.ID)
		if tokenErr != nil {
			return nil, tokenErr
		}
		user.AuthToken = authToken
		return user.Protobuf()
	}

	// Else signup
	if err := auth.DAL.CreateUser(&user); err != nil {
		return nil, err
	}
	authToken, tokenErr := auth.NewAuthToken(user.ID)
	if tokenErr != nil {
		return nil, tokenErr
	}
	user.AuthToken = authToken
	return user.Protobuf()
}

func (auth *Auth) getUserToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("Error getting the metadata of the context")
	}
	tokenSlice := md[auth.Config.AuthTokenHeader]
	if len(tokenSlice) < 1 {
		return "", fmt.Errorf("Token header not set")
	}
	return tokenSlice[0], nil
}

func (auth *Auth) getUserID(ctx context.Context) (string, error) {
	token, err := auth.getUserToken(ctx)
	if err != nil {
		return "", err
	}
	jwt := helpers.JWTHelper{HashSecretBytes: auth.Config.HashSecretBytes, Token: token}
	jwtTokenResult := jwt.Validate(auth.Config.JwtClaimUserID)
	switch jwtTokenResult.Status {
	case helpers.JWTokenStatusValid:
		if _, err := uuid.Parse(jwtTokenResult.Value); err != nil {
			return "", err
		}
		return jwtTokenResult.Value, nil
	case helpers.JWTokenStatusExpired:
		return "", fmt.Errorf("JWT token is expired")
	case helpers.JWTokenStatusInvalid, helpers.JWTokenNotAvailableYet:
		return "", fmt.Errorf("JWT token is not valid")
	}
	return "", fmt.Errorf("Unexpected status of the JWT token")
}

// NewAuthToken creates a new auth token for a given user id
func (auth *Auth) NewAuthToken(ID string) (string, error) {
	claims := make(map[string]interface{}, 2)
	claims[auth.Config.JwtClaimUserID] = ID
	jwt := helpers.JWTHelper{HashSecretBytes: auth.Config.HashSecretBytes}
	err := jwt.Generate(claims, auth.Config.JwtUserTokenDuration)
	return jwt.Token, err
}
