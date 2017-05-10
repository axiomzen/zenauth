package analytics

import (
// "github.com/sjhitchner/go-mixpanel" (old)
// "github.com/timehop/go-mixpanel" (newer, but needs https and a batch requests merged in and only 2 tests)
// https://github.com/dukex/mixpanel (more tests but less activity)
//
)

// Mixpanel the wrapper struct
type Mixpanel struct {
	// TODO: decide on mixpanel client library (or write our own)
	//api *mixpanel.MixpanelClient
}

// SendEvent example analytics event
func (a *Mixpanel) SendEvent(event *Event) error {
	// TODO: decide what is better to do; have an independent variable be true or false
	// for enabled, or calculate enabled such that its only on certain environments
	// if !(constants.ENV == constants.PRODUCTION || constants.ENV == constants.STAGING || constants.ENV == constants.CURATION) {
	// 	// Do not send if we are not on production/curation/staging
	// 	return nil
	// }

	// eventBody := bytes.NewReader(event.ToByteArray())
	// if mxpReq, err := http.NewRequest("POST", constants.ANALYTICS_EVENT_TRACK_URL, eventBody); err != nil {
	// 	return err
	// } else {
	// 	mxpReq.Header.Add("x-api-token", constants.API_KEY)
	// 	mxpReq.Header.Add("Content-Type", "application/json")
	// 	mxpReq.Close = true
	// 	if resp, err := http.DefaultClient.Do(mxpReq); err != nil {
	// 		log.Error("Could not send MixPanel event to analytics middleware")
	// 		return err
	// 	} else {
	// 		defer resp.Body.Close()
	// 		if resp.StatusCode != http.StatusOK {
	// 			return NewError("Mix Panel Event not OK, Status: " + resp.Status)
	// 		}
	// 	}
	// }
	return nil

}
