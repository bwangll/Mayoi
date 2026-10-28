// Command icongen places Mayoi's app icon art, cmd/mayoi/icon-art.png, on
// the macOS icon grid and writes resources/icon.png, the icon of builds
// that cannot compile cmd/mayoi/AppIcon.icon:
//
//	go run ./tools/icongen
package main

import (
	"image/png"
	"log"
	"os"

	"github.com/egoist/mygo/ui"
)

func main() {
	f, err := os.Open("cmd/mayoi/icon-art.png")
	if err != nil {
		log.Fatal(err)
	}
	art, err := png.Decode(f)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	bitmap := ui.NewBitmap(art)
	view := func(c *ui.Context) {
		c.Root().Background(ui.Transparent)
		ui.Column(c).Fill().Center().Children(func() {
			// The macOS icon grid: an 824 DIP squircle on a 1024 canvas.
			// macOS 26 shrinks an icon that does not fit it onto a gray frame.
			ui.Column(c).Size(824, 824).Radius(185).Clip().
				Shadow(0, 18, 40, 0, ui.RGBA(0, 0, 0, 0.22)).
				Children(func() {
					ui.Image(c, bitmap).Size(824, 824)
				})
		})
	}
	img := ui.Render(view, 1024, 1024, 1)
	out, err := os.Create("resources/icon.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// An *image.RGBA encodes as an 8-bit PNG; macOS 26 frames an icns of
	// 16-bit ones as a legacy icon.
	if err := png.Encode(out, img); err != nil {
		log.Fatal(err)
	}
}
