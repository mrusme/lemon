package unicorn

import "github.com/mrusme/lemon/inbox"

type Unicorn struct{}

func (out *Unicorn) Setup() error {
	return nil
}

func (out *Unicorn) Cleanup() {
}

func (out *Unicorn) Display(ibxMsg *inbox.Message) error {
	return nil
}
