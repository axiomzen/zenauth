// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package analytics

//"github.com/axiomzen/zenauth/models"

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

// ZENAUTHAnalyticsProvider is the data provider for this app
// TODO: consistency in API usage (pass in struct or return it?)
type ZENAUTHAnalyticsProvider interface {
	SendEvent(event *Event) error
}
