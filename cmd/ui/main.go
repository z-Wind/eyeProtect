package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const fontSize = 48

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
	s := fmt.Sprintf("%s\n休息倒數 %02d 秒", remindText, g.counter)
	f := &text.GoTextFace{Source: mplusSource, Size: fontSize}

	// 定義統一的行間距（以像素為單位）
	lineSpacing := fontSize * 1.5

	// 測量時必須傳入正確的 lineSpacing
	tw, th := text.Measure(s, f, lineSpacing)
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()

	op := &text.DrawOptions{}
	op.LineSpacing = lineSpacing
	// 讓文字區塊整體居中
	op.GeoM.Translate((float64(sw)-tw)/2, (float64(sh)-th)/2)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, s, f, op)
}

func (g *Game) Layout(w, h int) (int, int) { return w, h }

func main() {
	flag.IntVar(&waitSeconds, "w", 20, "倒數秒數")
	flag.BoolVar(&topEnable, "t", false, "置頂模式")
	flag.StringVar(&remindText, "r", "眼睛休息一下吧", "提醒文字")
	flag.Parse()

	var source *text.GoTextFaceSource
	fontPaths := []string{
		`C:\Windows\Fonts\msjh.ttc`,
		`C:\Windows\Fonts\msjh.ttf`,
		`/usr/share/fonts/truetype/wqy/wqy-microhei.ttc`,
	}

	for _, path := range fontPaths {
		fontData, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		// 關鍵修正：針對 .ttc 結尾的檔案使用 Collection 解析器
		if strings.HasSuffix(strings.ToLower(path), ".ttc") {
			sources, err := text.NewGoTextFaceSourcesFromCollection(bytes.NewReader(fontData))
			if err == nil && len(sources) > 0 {
				source = sources[0] // 取得合集中的第一個字體（通常是常規體）
				fmt.Printf("成功從合集載入字體: %s\n", path)
				break
			}
		} else {
			// 一般 .ttf 檔案
			s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
			if err == nil {
				source = s
				fmt.Printf("成功載入字體: %s\n", path)
				break
			}
		}
	}

	if source == nil {
		fmt.Println("注意：無法解析系統中文字體，改用內建字體")
		source, _ = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	}
	mplusSource = source

	ebiten.SetFullscreen(true)
	if topEnable {
		ebiten.SetWindowFloating(true)
	}
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&Game{counter: waitSeconds, timer: time.Now()}); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
