package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	var (
		intervalMin int
		waitSec     int
		topEnable   bool
		remindText  string
	)

	flag.IntVar(&intervalMin, "i", 10, "每隔多久提醒 (分鐘)")
	flag.IntVar(&waitSec, "w", 20, "休息時長 (秒)")
	flag.BoolVar(&topEnable, "t", true, "開啟置頂")
	flag.StringVar(&remindText, "r", "喝口水，站起來動一動", "自定義文字")
	flag.Parse()

	// 取得 UI 程式的路徑
	exeName := "eyeProtect"
	if runtime.GOOS == "windows" {
		exeName = "eyeProtect.exe"
	}
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	uiPath := filepath.Join(dir, exeName)

	log.Printf("Daemon 啟動：每 %d 分鐘提醒一次\n", intervalMin)

	for {
		time.Sleep(time.Duration(intervalMin) * time.Minute)

		// 構建指令參數
		args := []string{
			"-w", fmt.Sprintf("%d", waitSec),
			"-r", remindText,
		}
		if topEnable {
			args = append(args, "-t")
		}

		log.Println("觸發護眼視窗...")
		cmd := exec.Command(uiPath, args...)

		// 阻塞直到 UI 關閉
		if err := cmd.Run(); err != nil {
			log.Printf("視窗非正常關閉: %v\n", err)
		}
	}
}
