package dummy

import "github.com/mrusme/lemon/inbox"

type Dummy struct {
}

func (src *Dummy) Setup(ibx chan inbox.Message, opts interface{}) error {
	return nil
}

func (src *Dummy) Cleanup() {
}

func (src *Dummy) Start() (int, error) {
	for {
	}

	return 0, nil
}
