package main

import (
	"fmt"
	"image/color"
	"time"

	"image"

	"github.com/hajimehoshi/ebiten/v2"
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

var (
	counter int
	run     bool
)

func initFunc(prevScene string, payload interface{}) {
	run = true
	counter = 20

	fg := render.FontGenerator{
		File:    "mplus-1p-regular.ttf",
		Color:   image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Size:    size,
		Hinting: "",
		DPI:     dpi,
	}

	font := fg.Generate()
	text := font.NewStrText(fmt.Sprintf("眼睛休息 %02d 秒", counter), 0, 0)
	im := text.ToSprite().Bounds()
	x := oak.ScreenWidth/2 - im.Dx()/2
	y := oak.ScreenHeight/2 - im.Dy()/2
	text.SetPos(float64(x), float64(y))

	r, _ := render.Draw(text, 0)

	start := time.Now()
	event.GlobalBind(func(_ int, frames interface{}) int {
		end := time.Now()
		if end.Sub(start) >= time.Second {
			start = end
			r.Undraw()

			text = font.NewStrText(fmt.Sprintf("眼睛休息 %02d 秒", counter), float64(x), float64(y))

			r, _ = render.Draw(text, 0)
			if counter < 0 {
				// restart
				run = false
			}
			counter--
		}
		return 0
	}, event.Enter)
}

func oakMain() {
	oak.SetupConfig.Title = "eyeProtect"
	oak.SetupConfig.Screen.Width = screenWidth
	oak.SetupConfig.Screen.Height = screenHeight
	oak.SetupConfig.Assets.AssetPath = "./"
	oak.SetFullScreen(true)

	if err := oak.SetFullScreen(true); err != nil {
		fmt.Printf("%v", err)
		w, h := ebiten.ScreenSizeInFullscreen()
		oak.SetupConfig.Screen.Width = w
		oak.SetupConfig.Screen.Height = h
	}

	oak.Add("eyeProtect",
		// Init
		initFunc,
		// Loop
		func() bool {
			return run
		},

		// End
		func() (string, *scene.Result) {
			return "eyeProtect", nil
		},
	)

	render.SetDrawStack(
		render.NewHeap(false),
		// render.NewDrawFPS(),
	)
	oak.Init("eyeProtect")
}
