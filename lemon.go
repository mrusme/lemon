package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mrusme/lemon/inbox"
	"github.com/mrusme/lemon/output"
	"github.com/mrusme/lemon/pushover"
)

func main() {
	deviceId := os.Getenv("PUSHOVER_DEVICE_ID")
	secret := os.Getenv("PUSHOVER_SECRET")
	outputsString := os.Getenv("LEMON_OUTPUTS")

	var outputs []output.Output
	for _, outputString := range strings.Split(outputsString, ",") {
		o, err := output.New(outputString)
		if err != nil {
			panic(err)
		}
		outputs = append(outputs, o)
	}

	ibx := make(chan inbox.Message)

	po, err := pushover.New(ibx, deviceId, secret)
	if err != nil {
		panic(err)
	}

	go po.Stream()

	for {
		select {
		case ibxMsg := <-ibx:
			fmt.Println("Got new ibxMessage")
			fmt.Println(ibxMsg)

			for _, o := range outputs {
				if err := o.Display(&ibxMsg); err != nil {
					log.Printf("ERROR: %s\n", err)
				}
			}
		}
	}

	for _, o := range outputs {
		o.Cleanup()
	}

	os.Exit(0)
}
