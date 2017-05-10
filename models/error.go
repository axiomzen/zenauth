package models

import (
	"bytes"
	"runtime"
	"strconv"

	"github.com/axiomzen/authentication/constants"
)

//go:generate ffjson $GOFILE

// ErrorResponse is our api's returned error model
type ErrorResponse struct {
	Code         constants.APIErrorCode `form:"code"               json:"code"`
	ErrorMessage string                 `form:"error"              json:"error"`
	// boo; ffjson messes up if you use the error interface
	GoError AZError `form:"internalError"        json:"internalError"`
	//GoError string `form:"serverError"        json:"serverError"`
}

// AZError is our own underlying error representation
type AZError struct {
	GoError      string `form:"message"            json:"message"`
	CodeLocation string `form:"location"           json:"location"`
}

// Error satisfies the error interface
func (e *ErrorResponse) Error() string {
	return e.GoError.Error()
}

// Error statisfies the error interface
func (e *AZError) Error() string {
	if e == nil {
		return ""
	}
	return e.GoError
}

// Location is the line at which the erorr occurred
func (e *AZError) Location() string {
	return e.CodeLocation
}

// NewAZError creates a new underlying Error (where we are the ones generating the error)
func NewAZError(message string) *AZError {
	location := getErrorLocation()
	err := AZError{message, location}
	return &err
}

func getErrorLocation() string {
	pc, file, line, _ := runtime.Caller(2)

	var buffer bytes.Buffer

	buffer.WriteString(file)
	buffer.WriteString(": [")
	buffer.WriteString(runtime.FuncForPC(pc).Name())
	buffer.WriteString(" - line ")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString("]")
	return buffer.String()
}

// NewErrorResponse creates a new outer APIError
func NewErrorResponse(apiErrorCode constants.APIErrorCode, err *AZError, msgs ...string) *ErrorResponse {

	var buffer bytes.Buffer
	for _, m := range msgs {
		buffer.WriteString(m)
		buffer.WriteString(" ")
	}
	if buffer.Len() > 0 {
		buffer.Truncate(buffer.Len() - 1)
	}
	msg := buffer.String()

	if err == nil {
		// Use some Generic error message if the error is nil
		location := getErrorLocation()
		err = &AZError{"Error!", location}
	}

	apiErr := ErrorResponse{
		Code:         apiErrorCode,
		ErrorMessage: msg,
		GoError:      *err,
	}
	return &apiErr
}
