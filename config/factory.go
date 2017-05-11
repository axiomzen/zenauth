// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE NEED NOT BE EDITED BY HAND

package config

import (
	"sync"

	"github.com/axiomzen/envconfig"
)

// the config instance
var instance *ZENAUTHConfig
var once sync.Once
var initerr error

// Get retrieves the dispatcher config var
// returns an error if there is a problem
// re: http://marcio.io/2015/07/singleton-pattern-in-go/
func Get() (*ZENAUTHConfig, error) {
	once.Do(func() {
		var conf ZENAUTHConfig
		// prefix from hatch
		initerr = envconfig.Process("ZENAUTH", &conf)
		instance = &conf
		if initerr == nil {
			initerr = instance.ComputeDependents()
		}
	})
	return instance, initerr
}
