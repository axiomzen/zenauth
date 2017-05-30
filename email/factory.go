// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT

package email

import (
	"sync"

	"github.com/axiomzen/zenauth/config"
	"github.com/mailgun/mailgun-go"
)

var instance ZENAUTHEmailProvider
var once sync.Once
var initerr error

// Get retrieves the dispatcher config var
// returns an error if there is a problem
// re: http://marcio.io/2015/07/singleton-pattern-in-go/
func Get(conf *config.ZENAUTHConfig) (ZENAUTHEmailProvider, error) {
	once.Do(func() {
		if conf.EmailEnabled {
			instance = &MailgunImpl{
				gun: mailgun.NewMailgun(conf.MailGunDomain, conf.MailGunPrivateKey, conf.MailGunPublicKey),
			}
		} else {
			instance = &Noop{}
		}
	})
	return instance, initerr
}
