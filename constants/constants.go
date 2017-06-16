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
)
const (
	// APINetworkError for network errors
	APINetworkError = 5000 + iota
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
