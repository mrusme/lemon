package inbox

import "image"

type Priority uint8

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
)

type Message struct {
	Icon     image.Image
	IconPath string
	Title    string
	Text     string
	URL      string
	Prio     Priority
}
