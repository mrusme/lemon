package unicorn

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/mrusme/lemon/inbox"
	"golang.org/x/image/draw"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/unicornhd"
	"periph.io/x/host/v3"
)

type Unicorn struct {
	port spi.PortCloser
	uchd *unicornhd.Dev
}

func (out *Unicorn) Setup() error {
	var err error

	if _, err = host.Init(); err != nil {
		return err
	}

	out.port, err = spireg.Open("")
	if err != nil {
		return err
	}

	out.uchd, err = unicornhd.New(out.port)
	if err != nil {
		return err
	}

	out.uchd.Draw(out.uchd.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{})
	time.Sleep(time.Second * 1)
	out.uchd.Draw(out.uchd.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{})

	return nil
}

func (out *Unicorn) Cleanup() {
	if err := out.port.Close(); err != nil {
	}
}

func (out *Unicorn) Display(ibxMsg *inbox.Message) error {
	p := image.Pt(0, 0)
	icon := image.NewNRGBA(image.Rect(0, 0, out.uchd.Bounds().Dx(), out.uchd.Bounds().Dy()))
	draw.NearestNeighbor.Scale(icon, icon.Rect, ibxMsg.Icon, ibxMsg.Icon.Bounds(), draw.Over, nil)
	for i := 0; i < out.uchd.Bounds().Dx(); i++ {
		fmt.Println(i)
		p.X = -1 * i
		out.uchd.Draw(out.uchd.Bounds(), icon, icon.Bounds().Bounds().Min.Add(p))
		time.Sleep(time.Second * 1)
	}
	out.uchd.Draw(out.uchd.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{})
	return nil
}
