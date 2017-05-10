package data

import (
	"fmt"
)

// DALErrorCode the error code for the access layer errors
type DALErrorCode uint32

const (
	// DALErrorGeneral is for errors not specifically called out
	DALErrorGeneral DALErrorCode = iota
	// DALErrorCodeNoneAffected no rows were effected from the operation
	DALErrorCodeNoneAffected
	// DALErrorCodeUniqueEmail returned when the email already exists
	DALErrorCodeUniqueEmail
)

// DALError The error from the data access layer
type DALError struct {
	// Inner is the inner error (if present)
	Inner error
	// ErrorCode is the code for the error (so you don't have to look at strings)
	ErrorCode DALErrorCode
	text      string // errors that do not have a valid error just have text
}

// Error implements the error interface
func (e DALError) Error() string {
	if e.Inner != nil {
		return e.Inner.Error()
	} else if e.text != "" {
		return e.text
	}

	return fmt.Sprintf("Unknown Error: %d", e.ErrorCode)
}
