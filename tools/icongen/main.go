// Command icongen places Mayoi's app icon art, cmd/mayoi/icon-art.png, on
// the macOS icon grid and writes resources/icon.png:
//
//	go run ./tools/icongen
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/egoist/mygo/ui"
)

// frame is the width of the art's own pixel tile border (a dark outline
// and a light bevel), in pixels of the 1254 px art.
const frame = 52

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
	// macOS 26 shrinks an icon that does not fit its rounded tile onto a
	// gray frame, and the art's tile has square, stepped corners. So the
	// border is cropped off and the inside clipped to the tile instead.
	bitmap := ui.NewBitmap(glyph(art, art.Bounds().Inset(frame)))
	white := func(a float32) ui.Color { return ui.RGBA(255, 255, 255, a) }
	view := func(c *ui.Context) {
		c.Root().Background(ui.Transparent)
		ui.Column(c).Fill().Center().Children(func() {
			// The macOS icon grid: an 824 DIP squircle on a 1024 canvas.
			ui.Column(c).Size(824, 824).Radius(185).Clip().
				Shadow(0, 18, 40, 0, ui.RGBA(0, 0, 0, 0.22)).
				Gradient(white(0.55), white(0.18), 160).
				InsetShadow(0, 10, 24, 0, white(0.45)).
				Border(6, ui.Color{}).
				BorderLinearGradient(ui.LinearGradient{From: white(0.9), To: white(0.2), Angle: 160}).
				Children(func() {
					ui.Image(c, bitmap).Size(812, 812)
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

// glyph lifts the dark glyph off the art's flat tile color, the color at
// the art's corner, into a transparent image so it sits on the glass.
func glyph(art image.Image, r image.Rectangle) *image.NRGBA {
	tile := art.At(r.Min.X, r.Min.Y)
	tr, tg, tb, _ := tile.RGBA()
	out := image.NewNRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	ink := color.NRGBA{0x00, 0x2b, 0x36, 0}
	ir, ig, ib := float64(ink.R)*257, float64(ink.G)*257, float64(ink.B)*257
	span := math.Hypot(math.Hypot(float64(tr)-ir, float64(tg)-ig), float64(tb)-ib)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			pr, pg, pb, _ := art.At(x, y).RGBA()
			d := math.Hypot(math.Hypot(float64(pr)-float64(tr), float64(pg)-float64(tg)), float64(pb)-float64(tb))
			c := ink
			c.A = uint8(math.Round(255 * min(d/span, 1)))
			out.SetNRGBA(x-r.Min.X, y-r.Min.Y, c)
		}
	}
	return out
}
