package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/z-Wind/eyeProtect/internal/config"
)

// resolveUIPath 找出同目錄下的 UI 執行檔路徑
func resolveUIPath() (string, error) {
	exeName := "eyeProtect"
	if runtime.GOOS == "windows" {
		exeName = "eyeProtect.exe"
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", fmt.Errorf("無法解析執行檔目錄: %w", err)
	}
	uiPath := filepath.Join(dir, exeName)
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		return "", fmt.Errorf("找不到 UI 執行檔 %s，請確保它與 Daemon 放在一起", uiPath)
	}
	return uiPath, nil
}

// buildArgs 根據設定建構 UI 的命令列參數。
// 安全前提：使用 exec.Command 直接傳遞 []string，不經過 shell，
// 因此 RemindText 中的特殊字元不會造成命令注入風險。
func buildArgs(cfg config.Daemon) []string {
	args := []string{
		"-w", strconv.Itoa(cfg.WaitSec),
		"-r", cfg.RemindText,
	}
	if cfg.TopEnable {
		args = append(args, "-t")
	}
	return args
}

// triggerUI 啟動 UI 程式並阻塞直到視窗關閉。
// 接受 ctx 以支援 daemon 收到退出訊號時主動終止子程序。
func triggerUI(ctx context.Context, uiPath string, args []string) {
	log.Println("觸發護眼視窗...")
	cmd := exec.CommandContext(ctx, uiPath, args...)
	if err := cmd.Run(); err != nil {
		// exit status != 0（如使用者按 ESC 關閉，或 ctx 取消）屬正常情況
		log.Printf("視窗結束（可能為使用者手動關閉或收到退出訊號）: %v", err)
	}
	log.Println("休息結束，進入下一個計時循環。")
}

func main() {
	cfg := config.ParseDaemon()
	config.ValidateDaemon(cfg)

	uiPath, err := resolveUIPath()
	if err != nil {
		log.Fatalf("錯誤：%v", err)
	}

	log.Printf("Daemon 啟動：每 %d 分鐘提醒一次，休息 %d 秒", cfg.IntervalMin, cfg.WaitSec)

	// 監聽系統訊號，支援 Ctrl+C / SIGTERM 優雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	// 建立可取消的 context，用於在退出時中斷正在執行的 UI 子程序
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(time.Duration(cfg.IntervalMin) * time.Minute)
	defer ticker.Stop()

	args := buildArgs(cfg)

	for {
		select {
		case <-quit:
			log.Println("收到退出訊號，Daemon 停止。")
			cancel() // 主動終止正在執行的 UI 子程序
			return
		case <-ticker.C:
			triggerUI(ctx, uiPath, args)
		}
	}
}
