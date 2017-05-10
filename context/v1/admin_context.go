package v1

import (
	"github.com/axiomzen/authentication/helpers"
	"github.com/gocraft/web"

	"github.com/twinj/uuid"
)

// AdminContext for administrator authenticated routes
type AdminContext struct {
	*APIAuthContext

	AdminID string
}

// AuthRequired Middleware: Authorizes a user by authenticating the Json Web Token
func (c *AdminContext) AuthRequired(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes, Token: r.Header.Get(c.Config.AuthTokenHeader)}

	jwtTokenResult := jwt.Validate(c.Config.JwtClaimAdminID)

	switch jwtTokenResult.Status {
	case helpers.JWTokenStatusValid:

		if _, err := uuid.Parse(jwtTokenResult.Value); err != nil {
			c.UnauthorizedHandler(w, r)
			return
		}
		c.AdminID = jwtTokenResult.Value

		c.Log = c.Log.WithField("adminID", jwtTokenResult.Value)
		next(w, r)
	case helpers.JWTokenStatusExpired:
		c.ExpiredHandler(w, r)
	case helpers.JWTokenStatusInvalid, helpers.JWTokenNotAvailableYet:
		c.UnauthorizedHandler(w, r)
	}
}

// NewAuthToken creates a new auth token for a given admin id
func (c *AdminContext) NewAuthToken(ID string) (*helpers.JWTHelper, error) {
	claims := make(map[string]interface{}, 2)
	claims[c.Config.JwtClaimAdminID] = ID
	jwt := helpers.JWTHelper{HashSecretBytes: c.Config.HashSecretBytes}
	err := jwt.Generate(claims, c.Config.JwtAdminTokenDuration)
	return &jwt, err
}
