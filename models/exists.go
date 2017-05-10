package models

//go:generate ffjson $GOFILE

// Exists is a basic struct to return whether something exists
type Exists struct {
	Exists bool `json:"exists"`
}
