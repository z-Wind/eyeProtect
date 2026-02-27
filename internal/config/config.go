// Package config 定義 eyeProtect daemon 與 ui 兩個 binary 共用的參數常數與型別。
// 集中管理可避免兩端各自維護時產生預設值或欄位名稱不一致的問題。
package config

import (
	"flag"
	"log"
)

// 預設值常數，兩個 binary 統一引用
const (
	DefaultIntervalMin = 10
	DefaultWaitSec     = 20
	DefaultTopEnable   = true
	DefaultRemindText  = "喝口水，站起來動一動"

	DefaultUIRemindText = "眼睛休息一下吧"
)

// Daemon 集中管理 daemon 所有旗標參數
type Daemon struct {
	IntervalMin int
	WaitSec     int
	TopEnable   bool
	RemindText  string
}

// UI 集中管理 ui 所有旗標參數
type UI struct {
	WaitSeconds int
	TopEnable   bool
	RemindText  string
}

// ParseDaemon 解析並回傳 daemon 參數
func ParseDaemon() Daemon {
	var cfg Daemon
	flag.IntVar(&cfg.IntervalMin, "i", DefaultIntervalMin, "每隔多久提醒 (分鐘)")
	flag.IntVar(&cfg.WaitSec, "w", DefaultWaitSec, "休息時長 (秒)")
	flag.BoolVar(&cfg.TopEnable, "t", DefaultTopEnable, "開啟置頂")
	flag.StringVar(&cfg.RemindText, "r", DefaultRemindText, "自定義文字")
	flag.Parse()
	return cfg
}

// ParseUI 解析並回傳 ui 參數
func ParseUI() UI {
	var cfg UI
	flag.IntVar(&cfg.WaitSeconds, "w", DefaultWaitSec, "倒數秒數")
	flag.BoolVar(&cfg.TopEnable, "t", false, "置頂模式")
	flag.StringVar(&cfg.RemindText, "r", DefaultUIRemindText, "提醒文字")
	flag.Parse()
	return cfg
}

// ValidateDaemon 驗證 daemon 參數範圍，不合法時提早終止
func ValidateDaemon(cfg Daemon) {
	if cfg.IntervalMin < 1 {
		log.Fatalf("錯誤：-i 必須 >= 1（分鐘），目前值：%d", cfg.IntervalMin)
	}
	if cfg.WaitSec < 1 {
		log.Fatalf("錯誤：-w 必須 >= 1（秒），目前值：%d", cfg.WaitSec)
	}
}

// ValidateUI 驗證 ui 參數範圍，不合法時提早終止
func ValidateUI(cfg UI) {
	if cfg.WaitSeconds < 1 {
		log.Fatalf("錯誤：-w 必須 >= 1（秒），目前值：%d", cfg.WaitSeconds)
	}
}
