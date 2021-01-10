package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	mplusBigFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	counter int
	start   time.Time
}

func (g *Game) Update() error {
	end := time.Now()
	if end.Sub(g.start) >= time.Second {
		g.start = end
		g.counter--
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w, h := ebiten.WindowSize()
	s := fmt.Sprintf("眼睛休息 %02d 秒", g.counter)

	im := text.BoundString(mplusBigFont, s)
	text.Draw(screen, s, mplusBigFont, w/2-im.Dx()/2, h/2-im.Dy()/2, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func ebitenMain() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Eye Protect")
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(&Game{counter: 20, start: time.Now()}); err != nil {
		log.Fatal(err)
	}
}
