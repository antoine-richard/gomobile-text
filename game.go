// +build darwin linux

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	mfont "golang.org/x/mobile/exp/font"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

const (
	dpi = 72
)

type Game struct {
	font     *truetype.Font
	lastCalc clock.Time // when we last calculated a frame
}

func NewGame() *Game {
	var g Game
	g.reset()
	return &g
}

func (g *Game) loadFont() {
	f, err := asset.Open("System San Francisco Display Regular.ttf")
	if err != nil {
		fmt.Printf("error opening font asset: %v\n", err)
		g.loadFallbackFont()
		return
	}
	defer f.Close()
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("error reading font: %v\n", err)
		g.loadFallbackFont()
		return
	}
	g.font, err = freetype.ParseFont(raw)
	if err != nil {
		fmt.Printf("error parsing font: %v\n", err)
		g.loadFallbackFont()
		return
	}
}

func (g *Game) loadFallbackFont() {
	var err error
	fmt.Println("using Monospace font") // Default font doesn't work on Darwin
	g.font, err = truetype.Parse(mfont.Monospace())
	if err != nil {
		log.Fatalf("error parsing Monospace font: %v", err)
	}
}

func (g *Game) reset() {
	g.loadFont()
}

func (g *Game) Touch(down bool) {
	if down {
		fmt.Println("touch")
	}
}

func (g *Game) Update(now clock.Time) {
	// Compute game states up to now.
	for ; g.lastCalc < now; g.lastCalc++ {
		g.calcFrame()
	}
}

func (g *Game) calcFrame() {

}

type textSprite struct {
	text            string
	width           int
	height          int
	textColor       *image.Uniform
	backgroundColor *image.Uniform
	fontSize        float64
	x               geom.Pt
	y               geom.Pt
	leftAligned     bool
}

func (ts textSprite) draw(sz size.Event, g *Game, dynamicText string) {

	sprite := images.NewImage(ts.width, ts.height)

	// Background
	draw.Draw(sprite.RGBA, sprite.RGBA.Bounds(), ts.backgroundColor, image.ZP, draw.Src)

	d := &font.Drawer{
		Dst: sprite.RGBA,
		Src: ts.textColor,
		Face: truetype.NewFace(g.font, &truetype.Options{
			Size:    ts.fontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}

	dy := int(math.Ceil(ts.fontSize * dpi / dpi))
	textWidth := d.MeasureString(ts.text)

	// TODO: API to improve...
	if ts.leftAligned {
		d.Dot = fixed.Point26_6{
			X: fixed.I(0),
			Y: fixed.I(ts.height/2 + dy/2),
		}
	} else { // centered
		d.Dot = fixed.Point26_6{
			X: fixed.I(sz.Size().X/2) - (textWidth / 2),
			Y: fixed.I(ts.height/2 + dy/2),
		}
	}

	// TODO: API to improve...
	if dynamicText == "" {
		d.DrawString(ts.text)
	} else {
		d.DrawString(dynamicText)
	}

	// Draw the sprite on the screen
	sprite.Upload()
	sprite.Draw(
		sz,
		geom.Point{X: ts.x, Y: ts.y},
		geom.Point{X: ts.x + sz.WidthPt, Y: ts.y},
		geom.Point{X: ts.x, Y: ts.y + sz.HeightPt},
		sz.Bounds())
	sprite.Release()

}

func (g *Game) Render(sz size.Event, glctx gl.Context, images *glutil.Images) {

	loading := &textSprite{
		text:            "Loading...",
		width:           sz.WidthPx,
		height:          400,
		textColor:       image.White,
		backgroundColor: image.NewUniform(color.RGBA{0x35, 0x67, 0x99, 0xFF}),
		fontSize:        96,
		x:               0,
		y:               0,
	}

	text := "Loading" + strings.Repeat(".", int(time.Now().Unix()%4))
	loading.draw(sz, g, text)

	resolution := &textSprite{
		text:            fmt.Sprintf("%vpx * %vpx", sz.WidthPx, sz.HeightPx),
		width:           sz.WidthPx,
		height:          100,
		textColor:       image.White,
		backgroundColor: image.NewUniform(color.RGBA{0x31, 0xA6, 0xA2, 0xFF}),
		fontSize:        24,
		x:               0,
		y:               140, // ???? Pt?
		leftAligned:     true,
	}
	resolution.draw(sz, g, "")

}
