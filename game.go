// +build darwin linux

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/size"
	mfont "golang.org/x/mobile/exp/font"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

const (
	dpi           = 72
	smallFontSize = 32
	bigFontSize   = 100
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

func (g *Game) reset() {
	var err error
	
	g.font, err = freetype.ParseFont(mfont.Default())
	if err != nil {
		fmt.Println("Unable to parse default font, trying monospace")
		g.font, err = freetype.ParseFont(mfont.Monospace())
		if err != nil {
			log.Fatalf("error parsing monospace font: %v", err)
		}
	}

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

func (g *Game) Render(sz size.Event, glctx gl.Context, images *glutil.Images) {
	
	height := 400

	foreground := image.White
	background := image.NewUniform(color.RGBA{0x35, 0x67, 0x99, 0xFF})

	// Sprite to write text on
	textSprite := images.NewImage(sz.WidthPx, height/*sz.HeightPx*/)

	// Background to draw on
	draw.Draw(textSprite.RGBA, textSprite.RGBA.Bounds(), background, image.ZP, draw.Src)

	// Write the Loading... text on the sprite
	
	loadingText := "Loading" + strings.Repeat(".", int(time.Now().Unix()%4))
	
	d := &font.Drawer{
		Dst: textSprite.RGBA,
		Src: foreground,
		Face: truetype.NewFace(g.font, &truetype.Options{
			Size:    bigFontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	dy := int(math.Ceil(bigFontSize * dpi / 72))
	textWidth := d.MeasureString("Loading...")
	d.Dot = fixed.Point26_6{
		X: fixed.I(sz.Size().X/2) - (textWidth / 2),
		Y: fixed.I(height/*sz.Size().Y*//2 + dy/2),
	}
	d.DrawString(loadingText)
	
	// Write the resolution on the sprite
	
	resolutionText := fmt.Sprintf("%vpx * %vpx", sz.WidthPx, sz.HeightPx)
	d = &font.Drawer{
		Dst: textSprite.RGBA,
		Src: foreground,
		Face: truetype.NewFace(g.font, &truetype.Options{
			Size:    smallFontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	dy = int(math.Ceil(smallFontSize * dpi / 72))
	d.Dot = fixed.Point26_6{
		X: fixed.I(0),
		Y: fixed.I(height - dy/2),
	}
	d.DrawString(resolutionText)

	// Draw the text sprite on the screen
	textSprite.Upload()
	textSprite.Draw(
		sz,
		geom.Point{},
		geom.Point{X: sz.WidthPt},
		geom.Point{Y: sz.HeightPt},
		sz.Bounds())
	textSprite.Release()
	
}
