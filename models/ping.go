package models

//go:generate ffjson $GOFILE

// Ping is our basic health check response
type Ping struct {
	Ping string `form:"ping"                             json:"ping" lorem:",pong"`
}
