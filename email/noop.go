package email

// Noop does nothing
type Noop struct{}

// Send does nothing
func (n *Noop) Send(email *Message) error {
	return nil
}
