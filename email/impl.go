package email

import (
	//"github.com/mattbaird/gochimp"
	"bytes"
	"github.com/mailgun/mailgun-go"
	"io/ioutil"
)

// MailgunImpl the wrapper struct
type MailgunImpl struct {
	//api *gochimp.MandrillAPI
	gun mailgun.Mailgun
}

// Send sends an email
func (m *MailgunImpl) Send(email *Message) error {

	msg := m.gun.NewMessage(email.From, email.Subject, email.Body, email.To...)

	if len(email.AttachmentBody) > 0 {
		buf := &bytes.Buffer{}
		buf.Write(email.AttachmentBody)
		rc := ioutil.NopCloser(buf)
		msg.AddReaderAttachment(email.AttachmentFilename, rc)
	}

	// response, id, err
	_, _, err := m.gun.Send(msg)
	// TODO: perhaps log the ID of the message sent?
	// error check response?
	return err
}
