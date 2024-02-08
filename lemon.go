package main

import (
	"flag"
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

var flagSourcesString string
var flagOutputsString string

var flagPushoverDeviceID string
var flagPushoverSecret string

func env(name string, dflt string) string {
	val, exists := os.LookupEnv(name)
	if exists == false {
		return dflt
	}

	return val
}

func init() {
	flag.StringVar(
		&flagSourcesString,
		"sources",
		env("LEMON_SOURCES", "dummy"),
		"Notification sources to load, comma separated.\nAvailable: dummy pushover\nOverrides env LEMON_OUTPUTS\n")
	flag.StringVar(
		&flagOutputsString,
		"outputs",
		env("LEMON_OUTPUTS", "dbus"),
		"Notification outputs to load, comma separated.\nAvailable: dbus unicorn\nOverrides env LEMON_OUTPUTS\n")

	flag.StringVar(
		&flagPushoverDeviceID,
		"pushover-device-id",
		env("PUSHOVER_DEVICE_ID", ""),
		"Pushover source: The device ID to use.\nOverrides env PUSHOVER_DEVICE_ID\n")
	flag.StringVar(
		&flagPushoverSecret,
		"pushover-secret",
		env("PUSHOVER_SECRET", ""),
		"Pushover source: The secret to use.\nOverrides env PUSHOVER_SECRET\n")

}

func main() {
	flag.Parse()

	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, os.Interrupt, syscall.SIGTERM)

	ibx := make(chan inbox.Message)

	var sources []source.Source
	for _, sourceString := range strings.Split(flagSourcesString, ",") {
		var opts interface{}

		switch sourceString {
		case "pushover":
			opts = &pushover.PushoverOptions{
				DeviceID: flagPushoverDeviceID,
				Secret:   flagPushoverSecret,
			}
		}

		s, err := source.New(sourceString, ibx, opts)
		if err != nil {
			panic(err)
		}
		sources = append(sources, s)
	}

	var outputs []output.Output
	for _, outputString := range strings.Split(flagOutputsString, ",") {
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
