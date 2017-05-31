package email

import log "github.com/Sirupsen/logrus"

// Noop does nothing
type Noop struct{}

// Send does nothing
func (n *Noop) Send(email *Message) error {
	log.Debugf("Sending email:\n== HTML ==\n%v\n\n== TEXT ==\n%v", email.BodyHTML, email.Body)
	return nil
}
