package grpc

import (
	"fmt"

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
		linkToUser.Email = invite.GetInviteCode()
		linkUserErr = auth.DAL.GetUserByEmail(&linkToUser)
	case constants.InvitationTypeFacebook:
		linkToUser.FacebookID = invite.GetInviteCode()
		linkUserErr = auth.DAL.GetUserByFacebookID(&linkToUser)
	default:
		// Should never get here as we check the invite type above,
		linkUserErr = fmt.Errorf("Invitation type %s not supported", invite.GetType())
	}
	// Merge with calling user
	(&user).Merge(&linkToUser)

	var returnUser *protobuf.UserPublic
	var returnErr error
	if linkUserErr == nil {
		// User found, delete and return
		delUserErr := auth.DAL.DeleteUser(&linkToUser)
		auth.Log.WithError(delUserErr).Debug("Could not delete user while linking")
		returnUser, returnErr = linkToUser.ProtobufPublic()
		returnUser.Status = protobuf.UserStatus_merged
	} else {
		auth.Log.WithError(linkUserErr).Debug("Could not retrieve social account to link")
		returnUser, returnErr = user.ProtobufPublic()
	}
	if err := auth.DAL.UpdateUser(&user, &user); err != nil {
		return nil, err
	}

	return returnUser, returnErr
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
