// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE NEED NOT BE EDITED BY HAND

package config

import (
	"github.com/axiomzen/envconfig"
	"sync"
)

// the config instance
var instance *AUTHENTICATIONConfig
var once sync.Once
var initerr error

// Get retrieves the dispatcher config var
// returns an error if there is a problem
// re: http://marcio.io/2015/07/singleton-pattern-in-go/
func Get() (*AUTHENTICATIONConfig, error) {
	once.Do(func() {
		var conf AUTHENTICATIONConfig
		// prefix from hatch
		initerr = envconfig.Process("AXIOMZEN_AUTHENTICATION", &conf)
		instance = &conf
		if initerr == nil {
			initerr = instance.ComputeDependents()
		}
	})
	return instance, initerr
}
