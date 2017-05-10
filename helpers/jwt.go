// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND BUT IDEALLY YOU SHOULDN'T

package helpers

import (
	"fmt"
	"time"

	"errors"

	"github.com/twinj/uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// RE: https://stormpath.com/blog/jwt-the-right-way
// RE: https://littlemaninmyhead.wordpress.com/2015/11/22/cautionary-note-uuids-should-generally-not-be-used-for-authentication-tokens/

// JWTokenStatus is the status of the JWT
type JWTokenStatus int

const (
	// JWTokenStatusValid a valid token
	JWTokenStatusValid JWTokenStatus = iota
	// JWTokenStatusExpired means the token is expired
	JWTokenStatusExpired
	// JWTokenNotAvailableYet means the token isn't valid yet
	JWTokenNotAvailableYet
	// JWTokenStatusInvalid means the token in invalid
	JWTokenStatusInvalid

	exp string = "exp"
	iat string = "iat"
	jti string = "jti"
)

// JWTHelper is a nice wrapper for you
type JWTHelper struct {
	HashSecretBytes []byte
	Token           string
	JTI             string
}

// Generate creates a new jwt returns the erorr if there was any
func (h *JWTHelper) Generate(claims map[string]interface{}, d time.Duration) error {
	now := time.Now().UTC()
	//fmt.Printf("duration: %d\n", d)
	later := now.Add(d)
	claims[exp] = later.Unix()
	//fmt.Println("expiry time: " + later.String())
	// replay attacks prevention
	claims[iat] = now.Unix()
	// v4 is crypto-secure
	// allows for individual token invalidation (if we tie it back to the user in some way)
	uuidStr := uuid.NewV4().String()
	claims[jti] = uuidStr
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = jwt.MapClaims(claims)
	// TODO: Hide the actual error ?
	//fmt.Println("Generate hash secret bytes: " + string(h.HashSecretBytes))
	signedToken, err := token.SignedString(h.HashSecretBytes)

	if err != nil {
		return err
	}

	h.Token = signedToken
	h.JTI = uuidStr
	return nil
}

// JWTokenValidateResult the result struct
type JWTokenValidateResult struct {
	Message string
	Status  JWTokenStatus
	Value   string
}

// Validate checks the JWT Auth token, returning the status and
// the claim value if successful, "" otherwise
func (h *JWTHelper) Validate(claim string) *JWTokenValidateResult {
	// parse token
	token, err := jwt.Parse(h.Token, h.validateToken)
	// return if err
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				return &JWTokenValidateResult{Message: "Token expired", Status: JWTokenStatusExpired, Value: ""}
			}

			if ve.Errors&(jwt.ValidationErrorNotValidYet) != 0 {
				return &JWTokenValidateResult{Message: err.Error(), Status: JWTokenNotAvailableYet, Value: ""}
			}
			return &JWTokenValidateResult{Message: err.Error(), Status: JWTokenStatusInvalid, Value: ""}
		}

		return &JWTokenValidateResult{Message: err.Error(), Status: JWTokenStatusInvalid, Value: ""}
	}

	// return if token nil
	if token == nil {
		return &JWTokenValidateResult{Message: "Token was nil", Status: JWTokenStatusInvalid, Value: ""}
	}
	// check if token valid
	if !token.Valid {
		return &JWTokenValidateResult{Message: "Invalid token [0]", Status: JWTokenStatusInvalid, Value: ""}
	}
	// check expiry (done alrady by jwt)
	claims := token.Claims.(jwt.MapClaims)
	// exp, timeok := claims[exp].(float64)
	// // if no time, fail
	// if !timeok {
	// 	return &JWTokenValidateResult{Message: "Invalid token [1]", Status: JWTokenStatusInvalid, Value: ""}
	// }
	// now := time.Now().Unix()
	// if now > int64(exp) {
	// 	return &JWTokenValidateResult{Message: "Token expired", Status: JWTokenStatusExpired, Value: ""}
	// }
	// check for claim
	value, success := claims[claim].(string)

	if !success {
		return &JWTokenValidateResult{Message: "Invalid token [2]", Status: JWTokenStatusInvalid, Value: ""}
	}

	return &JWTokenValidateResult{Message: "Token Valid", Status: JWTokenStatusValid, Value: value}
}

// validateToken validates that the token is signed using the correct algorithm
// if so, it returns the secret, otherwise an error
// must have the signature: func(*Token) (interface{}, error)
func (h *JWTHelper) validateToken(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//fmt.Println("SIGNING METHOD DIFFERENT")
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return h.HashSecretBytes, nil
}

// Expire used for manually expiring a token (typically in testing scenarios)
func (h *JWTHelper) Expire() error {
	// parse token
	//fmt.Println("Expire Token: " + h.Token)
	//fmt.Println("Expire HashSecretBytes: " + string(h.HashSecretBytes))
	token, err := jwt.Parse(h.Token, h.validateToken)
	if err != nil {
		//fmt.Println("ERROR: " + err.Error())
		return errors.New("could not expire token [0]")
	}
	claims := token.Claims.(jwt.MapClaims)
	claims[exp] = time.Now().Unix()

	h.Token, err = token.SignedString(h.HashSecretBytes)
	if err != nil {
		return errors.New("could not expire token [1]")
	}
	return nil
}
