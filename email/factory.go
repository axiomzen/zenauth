// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT

package email

import (
	"github.com/axiomzen/authentication/config"
	"github.com/mailgun/mailgun-go"
	"sync"
)

var instance AUTHENTICATIONEmailProvider
var once sync.Once
var initerr error

// Get retrieves the dispatcher config var
// returns an error if there is a problem
// re: http://marcio.io/2015/07/singleton-pattern-in-go/
func Get(conf *config.AUTHENTICATIONConfig) (AUTHENTICATIONEmailProvider, error) {
	once.Do(func() {
		if conf.EmailEnabled {
			instance = &MailgunImpl{mailgun.NewMailgun(conf.MailGunDomain, conf.MailGunPrivateKey, conf.MailGunPublicKey)}
		} else {
			instance = &Noop{}
		}
	})
	return instance, initerr
}
