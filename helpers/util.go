// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package helpers

import (
	"strings"

	"github.com/axiomzen/null"
)

// IsZeroString returns true if the argument is null or
// the string pointed to by the argument is empty
func IsZeroString(s *string) bool {
	return s == nil || *s == ""
}

// IsZeroNullableString useful for when the value
// needs to be check for zero
func IsZeroNullableString(n null.String) bool {
	return !n.Valid || n.String == ""
}

// NewStringPointer returns a pointer to a copy of the string
// that is passed in
func NewStringPointer(str string) *string {
	//newstr := (str + " ")[:len(str)]
	newstr := str
	return &newstr
}

// EmailSanitize makes sure emails are space trimmed and lower cased
func EmailSanitize(email string) (safeEmail string) {
	safeEmail = strings.ToLower(email)
	safeEmail = strings.Trim(safeEmail, " \t")
	return
}
