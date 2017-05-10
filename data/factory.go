// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// YOU PROBABLY DO NOT NEED TO EDIT THIS FILE

package data

import (
	"errors"
	"github.com/axiomzen/authentication/config"
	"sync"

	"gopkg.in/pg.v4"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

// CreateProvider creates a new data provider
// Exposed for the tests - call Get() instead
func CreateProvider(conf *config.AUTHENTICATIONConfig) (AUTHENTICATIONProvider, error) {
	var err = errors.New("temp")
	var numtries uint16

	var provider AUTHENTICATIONProvider
	for err != nil && numtries < conf.PostgreSQLRetryNumTimes {
		pgdb := pg.Connect(&pg.Options{
			Addr:     conf.PostgreSQLHost + ":" + strconv.FormatUint(uint64(conf.PostgreSQLPort), 10),
			User:     conf.PostgreSQLUsername,
			Password: conf.PostgreSQLPassword,
			Database: conf.PostgreSQLDatabase,
			SSL:      *conf.PostgreSQLSSL,
		})
		provider = &dataProvider{db: pgdb}
		err = provider.Ping()
		if err != nil {
			log.WithFields(log.Fields{
				"numtries": numtries,
				"duration": conf.PostgreSQLRetrySleepTime,
				"port":     conf.PostgreSQLPort,
				"host":     conf.PostgreSQLHost,
				"database": conf.PostgreSQLDatabase,
			}).Info("Retrying connection...")
			// not sure if we need to close if we had an error
			provider.Close()
			// sleep
			time.Sleep(conf.PostgreSQLRetrySleepTime)
		}
		numtries++
	}
	return provider, err

}

var instance AUTHENTICATIONProvider
var once sync.Once
var initerr error

// Get retrieves the data provider instance
// returns an error if there is a problem
// re: http://marcio.io/2015/07/singleton-pattern-in-go/
func Get(conf *config.AUTHENTICATIONConfig) (AUTHENTICATIONProvider, error) {
	once.Do(func() {
		instance, initerr = CreateProvider(conf)
	})
	return instance, initerr
}
