// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package constants

import (
	"net/http"
	"time"
)

// HTTPStatusCode is the status codes that the app uses (a subset of all the available ones)
type HTTPStatusCode int

const (
	// StatusOK for GETs
	StatusOK HTTPStatusCode = http.StatusOK
	// StatusCreated for POSTs
	StatusCreated = http.StatusCreated
	// StatusNotFound not found
	StatusNotFound = http.StatusNotFound
	// StatusNoContent for PUTs
	StatusNoContent = http.StatusNoContent
	// StatusBadRequest malformed request
	StatusBadRequest = http.StatusBadRequest
	// StatusUnauthorized unauthorized (key or auth token)
	StatusUnauthorized = http.StatusUnauthorized
	// StatusForbidden unsure?
	StatusForbidden = http.StatusForbidden
	// StatusExpiredToken for expired tokens
	StatusExpiredToken = 440
	// StatusTokenNotAvailableYet for tokens that are not valid yet
	StatusTokenNotAvailableYet = 441
	// StatusServiceUnavailable not sure?
	StatusServiceUnavailable = http.StatusServiceUnavailable
	// StatusMovedPermanently for permantent redirects (http->https for example)
	StatusMovedPermanently = http.StatusMovedPermanently
	// StatusFound for forwarding requests
	StatusFound = http.StatusFound
	// StatusSeeOther for forwarding requests
	StatusSeeOther = http.StatusSeeOther
	// StatusTemporaryRedirect for posts?
	StatusTemporaryRedirect = http.StatusTemporaryRedirect
	// StatusInternalServerError for panics and other non handled errors
	StatusInternalServerError = http.StatusInternalServerError
)

//APIErrorCode The error codes that are sent back with errors (for debugging)
type APIErrorCode int

// New const brackets to reset iota
const (
	// APINotFound General not found error code
	APINotFound APIErrorCode = iota
	// APIInvalidRequest sent when we have a problem processing the request
	APIInvalidRequest
	// APIPanic sent when the api as recovered from a panic
	APIPanic
)

const (
	// APIDatabase a general error related to the database
	APIDatabase APIErrorCode = 1000 + iota
	// APIDatabaseUnreachable cannot reach the database
	APIDatabaseUnreachable
)

const (
	// APIDatabaseGet errors with getting data from the database
	APIDatabaseGet APIErrorCode = 1100 + iota
	// APIDatabaseGetUser error with retrieving a user
	APIDatabaseGetUser
)
const (
	// APIDatabaseCreate errors with inserting data
	APIDatabaseCreate APIErrorCode = 1200 + iota
	// APIDatabaseCreateUser errors creating users
	APIDatabaseCreateUser
	APIDatabaseCreateInvitation
)

const (
	// APIDatabaseUpdate errors with updating data
	APIDatabaseUpdate APIErrorCode = 1300 + iota
	// APIDatabaseUpdateUser errors updating users
	APIDatabaseUpdateUser
)
const (
	// APIDatabaseDelete errors with deleting data
	APIDatabaseDelete APIErrorCode = 1400 + iota
	// APIDatabaseDeleteUser deleting user
	APIDatabaseDeleteUser
)
const (
	// APIParsing Parsing
	APIParsing APIErrorCode = 2000 + iota
	// APIParsingUnsupportedContentType returned if we can't render that content type
	APIParsingUnsupportedContentType
	// APIParsingQueryParams returned if we have trouble parsing the query parameters
	APIParsingQueryParams
	// APIParsingMarshalling returned if we can't marshal our data structures
	APIParsingMarshalling
	// APIParsingJWT an error parsing the json web token
	APIParsingJWT
	// APIParsingUUIDUser an error occurred parsing the uuid string
	APIParsingUUIDUser
	// APIParsingPasswordHash has didn't work
	APIParsingPasswordHash
)
const (
	// APIGeneric generic errors
	APIGeneric APIErrorCode = 3000 + iota
)
const (
	// APIValidation request validation errors
	APIValidation APIErrorCode = 4000 + iota
	// APIValidationPasswordTooShort password too short
	APIValidationPasswordTooShort
	// APIValidationEmailNotValid email not valid
	APIValidationEmailNotValid
	APIValidationUserNameNotValid
)
const (
	// APINetworkError for network errors
	APINetworkError APIErrorCode = 5000 + iota
)
const (
	// APIUnauthorized api token incorrect
	APIUnauthorized APIErrorCode = 6000 + iota
	// APILoginSignupUnauthorized Login/Signup
	APILoginSignupUnauthorized
	// APILoginUserDoesNotExist user does not exist
	APILoginUserDoesNotExist
	// APIFacebookLoginNotValid facebook login invalid
	APIFacebookLoginNotValid
	// APIEmailInUse email is already in use
	APIEmailInUse
	// APISocialAccountExists social account already exists
	APISocialAccountExists
	// APILoginSignupInvalidCombination login/signup errors
	APILoginSignupInvalidCombination
	// APILoginNotVerified the email address has not been verified
	APILoginNotVerified
)

const (
	// APIAuthTokenCreation api token issue
	APIAuthTokenCreation APIErrorCode = 7000 + iota
	// APIExpiredAuthToken for expired tokens
	APIExpiredAuthToken
	// APIInvalidAuthToken for invalid tokens
	APIInvalidAuthToken
)

const (
	// APIForgotPasswordMessageError Error generating forget password emails
	APIForgotPasswordMessageError APIErrorCode = 8000 + iota
	// APIVerifyEmailMessageError Error generating email verification emails
	APIVerifyEmailMessageError
)

const (
	// APIInvitationsCreationError Error creating invitations
	APIInvitationsCreationError APIErrorCode = 9000 + iota
)

func (code APIErrorCode) String() string {
	switch code {
	case APINotFound:
		return "Not found"
	case APIInvalidRequest:
		return "Invalid request"
	case APIPanic:
		fallthrough
	case APIDatabase:
		fallthrough
	case APIDatabaseUnreachable:
		return "Server Error"
	case APIDatabaseGet:
		return "Could not get the requested resource"
	case APIDatabaseGetUser:
		return "Could not retrieve the user information"
	case APIDatabaseCreate:
		return "Could not create the resource"
	case APIDatabaseCreateUser:
		return "Could not create the user"
	case APIDatabaseCreateInvitation:
		return "Could not create the invitation"
	case APIDatabaseUpdate:
		return "Could not update the update"
	case APIDatabaseUpdateUser:
		return "Could not update your user info"
	case APIDatabaseDelete:
		return "Could not delete the resource"
	case APIDatabaseDeleteUser:
		return "Coul dnot delete the user"
	case APIParsing:
		fallthrough
	case APIParsingUnsupportedContentType:
		fallthrough
	case APIParsingQueryParams:
		fallthrough
	case APIParsingMarshalling:
		fallthrough
	case APIParsingJWT:
		fallthrough
	case APIParsingUUIDUser:
		fallthrough
	case APIParsingPasswordHash:
		return "Could not process your request"
	case APIGeneric:
		return "Could not process your request"
	case APIValidation:
		return ""
	case APIValidationPasswordTooShort:
		return "The password you provided was too short"
	case APIValidationEmailNotValid:
		return "The email was invalid"
	case APIValidationUserNameNotValid:
		return "The username was invalid"
	case APINetworkError:
		return ""
	case APIUnauthorized:
		return "You are not authorized to view this page"
	case APILoginSignupUnauthorized:
		return ""
	case APILoginUserDoesNotExist:
		return "User does not exist"
	case APIFacebookLoginNotValid:
		return "Could not login via facebook"
	case APIEmailInUse:
		return "Email already in use/exists"
	case APISocialAccountExists:
		return "Social account already exists"
	case APILoginSignupInvalidCombination:
		return "Invalid email/username/password combination"
	case APILoginNotVerified:
		return "User must validate their email first"
	case APIAuthTokenCreation:
		return "Could not login"
	case APIExpiredAuthToken:
		return "Session expired"
	case APIInvalidAuthToken:
		return "Session expired"
	case APIForgotPasswordMessageError:
		return "Unable to send password reset email"
	case APIVerifyEmailMessageError:
		return "Unable to send email verification"
	case APIInvitationsCreationError:
		return "Could not create invitations"
	default:
		return "Something went wrong!"
	}
}

// general constants
const (
	// EnvironmentStaging the staging environment
	EnvironmentStaging string = "STAGING"
	// EnvironmentProduction the production environment
	EnvironmentProduction string = "PRODUCTION"
	// EnvironmentTest the test environment (CI test usually)
	EnvironmentTest string = "TEST"
	// EnvironmentDevelopment local development environment
	EnvironmentDevelopment string = "DEVELOPMENT"
	// TimeFormat
	// This loses all nanosecond precision
	// TimeFormat = "2006-01-02T15:04:05Z0700"
	TimeFormat = time.RFC3339Nano

	InvitationTypeEmail    = "email"
	InvitationTypeFacebook = "facebook"
)

var (
	// InvitationTypes is used to check the types of invitations we handle
	InvitationTypes = map[string]bool{
		InvitationTypeEmail:    true,
		InvitationTypeFacebook: true,
	}
)

// in case we want to ever support multiple
// // DBType is the types of databases we support
// type DBType int

// const (
// 	// Mock is a fake database
// 	Mock DBType = iota
// 	//MongoDB the mongo DB database (Document)
// 	MongoDB
// 	// Postgres is the postgresql database (SQL)
// 	Postgres
// 	// CouchDB is the couch DB database (Document)
// 	CouchDB
// 	// CouchBase is the couchbase database (Document)
// 	CouchBase
// 	// OrientDb is the orientDb database (Hybrid Document or Graph)
// 	OrientDb
// 	// Redis is the redis database (Key-Value)
// 	Redis
// )
