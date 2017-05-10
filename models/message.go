package models

//go:generate ffjson $GOFILE

// Message is a simple message
type Message struct {
	Message string `json:"message" lorem:"sentence,2,10"`
}
