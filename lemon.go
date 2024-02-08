package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mrusme/lemon/inbox"
	"github.com/mrusme/lemon/output"
	"github.com/mrusme/lemon/source"
	"github.com/mrusme/lemon/source/pushover"
)

func main() {
	deviceId := os.Getenv("PUSHOVER_DEVICE_ID")
	secret := os.Getenv("PUSHOVER_SECRET")
	sourcesString := os.Getenv("LEMON_SOURCES")
	outputsString := os.Getenv("LEMON_OUTPUTS")

	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, os.Interrupt, syscall.SIGTERM)

	ibx := make(chan inbox.Message)

	var sources []source.Source
	for _, sourceString := range strings.Split(sourcesString, ",") {
		s, err := source.New(sourceString, ibx, &pushover.PushoverOptions{
			DeviceID: deviceId,
			Secret:   secret,
		})
		if err != nil {
			panic(err)
		}
		sources = append(sources, s)
	}

	var outputs []output.Output
	for _, outputString := range strings.Split(outputsString, ",") {
		o, err := output.New(outputString)
		if err != nil {
			panic(err)
		}
		outputs = append(outputs, o)
	}

	for _, source := range sources {
		go source.Start()
	}

	log.Println("All set, entering main loop ...")
mainloop:
	for {
		select {
		case ibxMsg := <-ibx:
			log.Println("Got new ibxMessage")
			log.Println(ibxMsg)

			for _, o := range outputs {
				if err := o.Display(&ibxMsg); err != nil {
					log.Printf("ERROR: %s\n", err)
				}
			}
		case sig := <-osSig:
			log.Printf("Received signal: %s\n", sig.String())
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				log.Println("Breaking main loop...")
				break mainloop
			}
		}
	}

	log.Println("Cleaning up ...")
	for _, o := range outputs {
		o.Cleanup()
	}

	log.Println("Bye!")
	os.Exit(0)
}
