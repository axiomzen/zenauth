// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT

package analytics

import (
	"github.com/axiomzen/authentication/config"
	"sync"
)

var instance AUTHENTICATIONAnalyticsProvider
var once sync.Once
var initerr error

// Get retrieves the dispatcher config var
// returns an error if there is a problem
// re: http://marcio.io/2015/07/singleton-pattern-in-go/
// TODO: perhaps we should be passed in the conf instead
func Get(conf *config.AUTHENTICATIONConfig) (AUTHENTICATIONAnalyticsProvider, error) {
	once.Do(func() {
		if conf.AnalyticsEnabled {
			// TODO: conf.MixPanelAPIKey
			instance = &Mixpanel{}
		} else {
			instance = &Noop{}
		}
	})
	return instance, initerr
}
