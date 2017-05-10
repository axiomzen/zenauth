// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package analytics

import (
//"github.com/axiomzen/authentication/models"
)

// Event example analytics event
type Event struct {
	EventName string  `json:"eventName"`
	Options   Options `json:"options"`
}

// Options example analytics options
type Options struct {
	DistinctID string `json:"distinct_id"`
	Email      string `json:"email,omitempty"`
	Created    string `json:"created,omitempty"`
	OS         string `json:"os,omitempty"`
}

// AUTHENTICATIONAnalyticsProvider is the data provider for this app
// TODO: consistency in API usage (pass in struct or return it?)
type AUTHENTICATIONAnalyticsProvider interface {
	SendEvent(event *Event) error
}
