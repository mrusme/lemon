package inbox

import "image"

type Message struct {
	Icon     image.Image
	IconPath string
	Title    string
	Text     string
	URL      string
}
