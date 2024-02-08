package dummy

import (
	"log"

	"github.com/mrusme/lemon/helpers"
	"github.com/mrusme/lemon/inbox"
)

type Dummy struct {
	ibx chan inbox.Message
}

func (src *Dummy) Setup(ibx chan inbox.Message, opts interface{}) error {
	src.ibx = ibx
	return nil
}

func (src *Dummy) Cleanup() {
}

func (src *Dummy) Start() (int, error) {
	for {
		icon, iconPath, err := helpers.GetIcon("", "dummy")
		if err != nil {
			log.Printf("Error: %s\n", err)
			continue
		}
		ibxMsg := inbox.Message{
			Icon:     icon,
			IconPath: iconPath,
			Title:    "Dummy Title ",
			Text:     "Dummy Message",
			URL:      "http://xn--gckvb8fzb.com",
		}
		src.ibx <- ibxMsg
	}

	return 0, nil
}
