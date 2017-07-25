// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE NEED NOT BE EDITED BY HAND

package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/axiomzen/envconfig"
	"github.com/axiomzen/zenauth/constants"
	"github.com/joho/godotenv"
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
		// set local dev env
		if strings.EqualFold(os.Getenv("ZENAUTH_ENVIRONMENT"), constants.EnvironmentDevelopment) {
			if err := godotenv.Load(); err != nil {
				// Do not fail here, incase they've manually loaded env variables
				// Will fail out at config.Get() if any require env variables not set
				fmt.Println("error loading godotenv")
			}
			fmt.Println(".env Loaded")
		}
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
