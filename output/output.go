package output

import (
	"errors"

	"github.com/mrusme/lemon/inbox"
	"github.com/mrusme/lemon/output/dbus"
	"github.com/mrusme/lemon/output/unicorn"
)

type Output interface {
	Setup() error
	Cleanup()
	Display(ibxMsg *inbox.Message) error
}

func New(name string) (Output, error) {
	var output Output

	switch name {
	case "dbus":
		output = new(dbus.Dbus)
	case "unicorn":
		output = new(unicorn.Unicorn)
	default:
		return nil, errors.New("No such output")
	}

	err := output.Setup()
	if err != nil {
		return nil, err
	}

	return output, nil
}
