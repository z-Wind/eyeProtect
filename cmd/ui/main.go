package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusSource *text.GoTextFaceSource
	waitSeconds int
	topEnable   bool
	remindText  string
)

type Game struct {
	counter int
	timer   time.Time
}

func (g *Game) Update() error {
	if g.counter < 0 || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if time.Since(g.timer) >= time.Second {
		g.timer = time.Now()
		g.counter--
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 顯示提醒文字與倒數
	s := fmt.Sprintf("%s\n休息倒數 %02d 秒", remindText, g.counter)
	f := &text.GoTextFace{Source: mplusSource, Size: 48}

	tw, th := text.Measure(s, f, 1.5)
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()

	op := &text.DrawOptions{}
	op.GeoM.Translate((float64(sw)-tw)/2, (float64(sh)-th)/2)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = th / 2 // 簡單的兩行間距
	text.Draw(screen, s, f, op)
}

func (g *Game) Layout(w, h int) (int, int) { return w, h }

func main() {
	flag.IntVar(&waitSeconds, "w", 20, "倒數秒數")
	flag.BoolVar(&topEnable, "t", false, "置頂模式")
	flag.StringVar(&remindText, "r", "眼睛休息一下吧", "提醒文字")
	flag.Parse()

	s, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	mplusSource = s

	ebiten.SetFullscreen(true)
	if topEnable {
		ebiten.SetWindowFloating(true)
	}
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&Game{counter: waitSeconds, timer: time.Now()}); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
