package source

import (
	"errors"

	"github.com/mrusme/lemon/inbox"
	"github.com/mrusme/lemon/source/pushover"
)

type Source interface {
	Setup(ibx chan inbox.Message, opts interface{}) error
	Cleanup()

	Start() (int, error)
}

func New(name string, ibx chan inbox.Message, opts interface{}) (Source, error) {
	var source Source

	switch name {
	case "pushover":
		source = new(pushover.Pushover)
	default:
		return nil, errors.New("No such source")
	}

	err := source.Setup(ibx, opts)
	if err != nil {
		return nil, err
	}

	return source, nil
}
