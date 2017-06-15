package grpc

import (
	"fmt"

	context "golang.org/x/net/context"

	"google.golang.org/grpc/metadata"

	"github.com/axiomzen/zenauth/config"
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
