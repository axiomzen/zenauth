// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package email

//"github.com/axiomzen/zenauth/models"

// Message is the basic holder for the info
type Message struct {
	From               string
	Subject            string
	Body               string
	BodyHTML           string
	To                 []string
	ReplyTo            string
	AttachmentFilename string
	AttachmentBody     []byte
}

// ZENAUTHEmailProvider is the data provider for this app
// TODO: where should the complexity live? if we want to swap out
// email providers, it would make sense to keep it out of here
// and keep the logic in the middleware (with helpers)
type ZENAUTHEmailProvider interface {
	// Send sends an email (what more do you want)
	Send(email *Message) error
}
