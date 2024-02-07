package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/unicornhd"
	"periph.io/x/host/v3"
)

var port spi.PortCloser

func setup() {
	var err error

	if _, err = host.Init(); err != nil {
		log.Fatal(err)
	}

	port, err = spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	uchd, err := unicornhd.New(port)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(uchd.String())

	uchd.Draw(uchd.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{})

	fmt.Println("Ok")
}

func bye() {

	if err := port.Close(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	defer bye()
	setup()

}
