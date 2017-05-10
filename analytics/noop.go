package analytics

// Noop does a nop for everything
type Noop struct{}

// SendEvent does a no-op
func (n *Noop) SendEvent(event *Event) error {
	return nil
}
