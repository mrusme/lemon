package unicorn

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"time"

	_ "embed"

	"github.com/golang/freetype/truetype"
	"github.com/mrusme/lemon/inbox"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/unicornhd"
	"periph.io/x/host/v3"
)

//go:embed "fonts/Hack-Regular.ttf"
var fs embed.FS

type Unicorn struct {
	port spi.PortCloser
	uchd *unicornhd.Dev

	fontTT     *truetype.Font
	fontFace   font.Face
	typewriter *font.Drawer
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

	fontBytes, err := fs.ReadFile("fonts/Hack-Regular.ttf")
	if err != nil {
		return err
	}

	out.fontTT, err = truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	out.fontFace = truetype.NewFace(out.fontTT, &truetype.Options{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	return nil
}

func (out *Unicorn) Cleanup() {
	if err := out.port.Close(); err != nil {
	}
}

func (out *Unicorn) Display(ibxMsg *inbox.Message) error {
	label := fmt.Sprintf("%s: %s", ibxMsg.Title, ibxMsg.Text)

	p := image.Pt(-out.uchd.Bounds().Dx(), 0)

	icon := image.NewNRGBA(image.Rect(0, 0, out.uchd.Bounds().Dx(), out.uchd.Bounds().Dy()))
	draw.NearestNeighbor.Scale(icon, icon.Rect, ibxMsg.Icon, ibxMsg.Icon.Bounds(), draw.Over, nil)

	// TODO: Fix quick & dirty hack of `len(label) * 12`
	tmpRect := image.Rectangle{image.Point{0, 0}, image.Point{len(label) * 12, 16}}
	tmp := image.NewNRGBA(tmpRect)
	draw.Draw(tmp, icon.Bounds(), icon, image.Point{0, 0}, draw.Src)

	out.typewriter = &font.Drawer{
		Dst:  tmp,
		Src:  &image.Uniform{color.NRGBA{255, 255, 255, 255}},
		Face: out.fontFace,
	}
	out.typewriter.Dot = fixed.Point26_6{fixed.I(22), fixed.I(12)}
	out.typewriter.DrawString(label)

	// DEBUG
	// outputFile, err := os.Create("combined_image.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer outputFile.Close()
	//
	// if err := png.Encode(outputFile, tmp); err != nil {
	// 	panic(err)
	// }
	// os.Exit(0)

	for i := 0; i < (out.uchd.Bounds().Dx() + tmp.Bounds().Dx()); i++ {
		p.X += 1
		out.uchd.Draw(out.uchd.Bounds(), tmp, tmp.Bounds().Bounds().Min.Add(p))
		time.Sleep(time.Millisecond * 100)
	}
	out.uchd.Draw(out.uchd.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{})
	return nil
}
