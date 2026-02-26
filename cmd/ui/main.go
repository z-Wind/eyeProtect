package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"

	// 注意：請務必維持使用 v1 text 包。
	// 理由：text/v2 在舊型整合顯卡（如 Intel HD 4000）上，其動態紋理快取（Atlas）
	// 與 Glyph 渲染引擎會觸發驅動程式 Bug，導致英數部分出現「細小方塊」或「黑塊」。
	// v1 版本採用較原始的紋理上傳方式，對舊硬體相容性最佳。
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const fontSize = 48

var (
	mplusFace   font.Face
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

// Draw 渲染邏輯說明：
// 1. 為了確保在舊顯卡（Intel HD 4000 等）穩定顯示，不使用 text/v2 的自動佈局。
// 2. 這裡採用手動拆分行（line1, line2）並計算座標。
// 3. 使用 text.Draw (v1) 配合固定的 font.Face，避免觸發顯卡的紋理記憶體對齊錯誤。
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 20, 255})

	line1 := remindText
	line2 := fmt.Sprintf("休息倒數 %02d 秒", g.counter)
	lines := []string{line1, line2}

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	lineHeight := int(fontSize * 1.5)
	totalHeight := len(lines) * lineHeight

	for i, line := range lines {
		rect := text.BoundString(mplusFace, line)
		tw := rect.Dx()
		x := (sw - tw) / 2
		// y 是 Baseline 座標，這樣算置中比較精準
		y := (sh-totalHeight)/2 + (i * lineHeight) + fontSize
		text.Draw(screen, line, mplusFace, x, y, color.White)
	}
}

func (g *Game) Layout(w, h int) (int, int) { return w, h }

func main() {
	// 【關鍵相容性設定】
	// 強制指定 OpenGL 繪圖庫。Intel HD 4000 在 Windows 上的 DX11/12 驅動極不穩定。
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")

	// 關閉外部 Shader 優化，防止舊驅動在處理動態文字紋理時發生位元偏移（Bit Shift）。
	os.Setenv("EBITEN_EXTERNAL_SHADER", "0")

	flag.IntVar(&waitSeconds, "w", 20, "倒數秒數")
	flag.BoolVar(&topEnable, "t", false, "置頂模式")
	flag.StringVar(&remindText, "r", "眼睛休息一下吧", "提醒文字")
	flag.Parse()

	// 載入字體數據
	var fontData []byte
	fontPath := `C:\Windows\Fonts\msjh.ttc`
	data, err := os.ReadFile(fontPath)
	if err != nil {
		fontData = fonts.MPlus1pRegular_ttf
	} else {
		fontData = data
	}

	// 修正後的解析邏輯
	collection, err := opentype.ParseCollection(fontData)
	var finalFont *opentype.Font
	if err != nil {
		// 如果不是 Collection，嘗試當作單一 TTF 解析
		finalFont, err = opentype.Parse(fontData)
		if err != nil {
			log.Fatal("無法解析字體數據:", err)
		}
	} else {
		// 從 Collection 中取得第一個字體
		finalFont, err = collection.Font(0)
		if err != nil {
			log.Fatal("無法從合集中取得字體:", err)
		}
	}

	// 建立 Face
	mplusFace, err = opentype.NewFace(finalFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull, // 舊顯卡建議開啟 Hinting 讓邊緣清晰
	})
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetFullscreen(true)
	if topEnable {
		ebiten.SetWindowFloating(true)
	}
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&Game{counter: waitSeconds, timer: time.Now()}); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
