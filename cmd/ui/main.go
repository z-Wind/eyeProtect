package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"

	// 注意：維持使用 v1 text 包。
	// text/v2 在舊型整合顯卡（如 Intel HD 4000）上，其動態紋理快取（Atlas）
	// 與 Glyph 渲染引擎會觸發驅動程式 Bug，導致英數部分出現「細小方塊」或「黑塊」。
	// v1 版本採用較原始的紋理上傳方式，對舊硬體相容性最佳。
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/z-Wind/eyeProtect/internal/config"
)

const (
	fontSize    = 48
	lineSpacing = 1.5 // 行距倍數
)

// loadFont 嘗試載入系統字體，所有路徑失敗時回退到內建字體。
// 永不回傳錯誤，因此簡化為單一回傳值。
func loadFont() []byte {
	candidates := []string{
		`C:\Windows\Fonts\msjh.ttc`,          // Windows 微軟正黑
		`C:\Windows\Fonts\mingliu.ttc`,       // Windows 細明體（備援）
		"/System/Library/Fonts/PingFang.ttc", // macOS
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err == nil {
			log.Printf("使用系統字體：%s", path)
			return data
		}
	}
	log.Println("找不到系統字體，使用內建 MPlus1p 字體")
	return fonts.MPlus1pRegular_ttf
}

// parseFontFace 從字體資料（TTC 或 TTF）解析出 font.Face
func parseFontFace(fontData []byte, dpi float64) (font.Face, error) {
	// 優先嘗試 Collection（.ttc）
	if col, err := opentype.ParseCollection(fontData); err == nil {
		f, err := col.Font(0)
		if err != nil {
			return nil, fmt.Errorf("無法從字體合集取得字體: %w", err)
		}
		return newFace(f, dpi)
	}
	// 退回單一 TTF
	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("無法解析字體資料: %w", err)
	}
	return newFace(f, dpi)
}

func newFace(f *opentype.Font, dpi float64) (font.Face, error) {
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull, // 舊顯卡建議開啟 Hinting 讓邊緣清晰
	})
}

// Game 實作 ebiten.Game 介面
type Game struct {
	face    font.Face
	lines   [2]string // 固定兩行：提醒文字 + 倒數
	counter int
	timer   time.Time
}

func newGame(cfg config.UI, face font.Face) *Game {
	g := &Game{
		face:    face,
		counter: cfg.WaitSeconds,
		timer:   time.Now(),
	}
	g.lines[0] = cfg.RemindText
	g.lines[1] = fmt.Sprintf("休息倒數 %02d 秒", cfg.WaitSeconds)
	return g
}

func (g *Game) Update() error {
	if g.counter <= 0 || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if time.Since(g.timer) >= time.Second {
		// 用 Add(time.Second) 推算下一個基準點，避免累積誤差
		g.timer = g.timer.Add(time.Second)
		g.counter--
		// 僅在秒數變化時重新格式化，每秒最多執行一次
		g.lines[1] = fmt.Sprintf("休息倒數 %02d 秒", g.counter)
	}
	return nil
}

// Draw 使用 text.Draw (v1) + 手動換行，確保舊顯卡相容性
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 20, 255})

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	lineHeight := int(fontSize * lineSpacing)
	totalHeight := len(g.lines) * lineHeight

	for i, line := range g.lines {
		rect := text.BoundString(g.face, line)
		x := (sw - rect.Dx()) / 2
		// clamp x，避免超長文字時座標為負數導致左側截斷
		if x < 0 {
			x = 0
		}
		y := (sh-totalHeight)/2 + i*lineHeight + fontSize
		text.Draw(screen, line, g.face, x, y, color.White)
	}
}

func (g *Game) Layout(w, h int) (int, int) { return w, h }

// setEnv 設定環境變數，失敗時記錄警告
func setEnv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Printf("警告：無法設定環境變數 %s: %v", key, err)
	}
}

func main() {
	// 【關鍵相容性設定】
	// 強制 OpenGL：Intel HD 4000 在 Windows 上的 DX11/12 驅動極不穩定。
	setEnv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	// 關閉外部 Shader 優化，防止舊驅動處理動態文字紋理時發生位元偏移。
	setEnv("EBITEN_EXTERNAL_SHADER", "0")

	cfg := config.ParseUI()
	config.ValidateUI(cfg)

	fontData := loadFont()

	// 透過 DeviceScaleFactor 動態計算 DPI，支援 HiDPI / Retina 螢幕
	scaleFactor := ebiten.Monitor().DeviceScaleFactor()
	dpi := 72.0 * scaleFactor

	face, err := parseFontFace(fontData, dpi)
	if err != nil {
		log.Fatal("建立字體 Face 失敗:", err)
	}
	defer face.Close() // 明確釋放字形快取資源

	ebiten.SetFullscreen(true)
	ebiten.SetWindowFloating(cfg.TopEnable)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	// 降低 TPS，純倒數畫面不需要 60 FPS，大幅減少 CPU/GPU 空轉
	ebiten.SetTPS(2)

	game := newGame(cfg, face)
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
