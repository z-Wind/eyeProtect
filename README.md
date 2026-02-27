# eyeProtect

依照 **20-20-20 護眼原則** 設計的輕量桌面提醒工具：每隔 N 分鐘自動彈出全螢幕視窗，倒數計時結束後自動關閉，提醒你暫時離開螢幕讓眼睛休息。

---

## 架構說明

專案由兩個獨立執行檔組成，職責分離：

```
eyeProtect/
├── cmd/
│   ├── daemon/   ← 背景計時主控端（daemon）
│   └── ui/       ← 全螢幕提醒視窗（ui）
└── internal/
    └── config/   ← 兩端共用的參數定義與驗證邏輯
```

| 執行檔 | 職責 |
|--------|------|
| `daemon` | 背景長駐，依照設定的間隔週期性呼叫 `eyeProtect`（ui） |
| `eyeProtect` | 顯示全螢幕深色倒數視窗，計時結束或按下 ESC 後自動退出 |

**通訊方式：** daemon 以命令列參數（`-w`、`-r`、`-t`）將設定傳遞給 ui，使用 `exec.CommandContext` 啟動並阻塞等待，daemon 收到退出訊號時會同時終止正在執行的 ui 子程序。

---

## 編譯方式

### 前置需求

- Go 1.21+
- Windows / macOS / Linux（需有圖形環境）

### Windows

```bat
REM 編譯 UI 顯示端
go build -ldflags="-H windowsgui" -o eyeProtect.exe ./cmd/ui

REM 編譯背景主控端
go build -o daemon.exe ./cmd/daemon
```

> `-H windowsgui` 可讓 daemon 啟動 UI 時不額外彈出黑色命令提示字元視窗。

### macOS / Linux

```bash
# 編譯 UI 顯示端
go build -o eyeProtect ./cmd/ui

# 編譯背景主控端
go build -o daemon ./cmd/daemon
```

將兩個執行檔放在**同一個目錄**下，daemon 會自動在相同目錄尋找 `eyeProtect`（或 `eyeProtect.exe`）。

---

## 參數說明

### daemon 參數

| 旗標 | 預設值 | 說明 |
|------|--------|------|
| `-i` | `10` | 每隔多久提醒一次（分鐘，最小值 1） |
| `-w` | `20` | 休息視窗顯示時長（秒，最小值 1） |
| `-t` | `true` | 護眼視窗是否置頂 |
| `-r` | `喝口水，站起來動一動` | 顯示在視窗上的提醒文字 |

### eyeProtect（ui）參數

ui 通常由 daemon 自動呼叫，也可單獨執行測試：

| 旗標 | 預設值 | 說明 |
|------|--------|------|
| `-w` | `20` | 倒數秒數（最小值 1） |
| `-t` | `false` | 視窗是否置頂 |
| `-r` | `眼睛休息一下吧` | 提醒文字 |

---

## 執行範例

```bash
# 每 25 分鐘提醒一次，休息 20 秒，自定義提醒文字
daemon -i 25 -w 20 -r "看看窗外，動動肩頸"

# 關閉置頂（方便在提醒時繼續操作其他視窗）
daemon -i 10 -w 30 -t=false

# 單獨測試 UI 視窗（顯示 5 秒後自動關閉）
./eyeProtect -w 5 -r "測試中"
```

---

## 開機自動啟動

### Windows — 工作排程器

1. 開啟「工作排程器」→「建立基本工作」
2. 觸發程序：**登入時**
3. 動作：啟動程式，填入 `daemon.exe` 的完整路徑與所需參數
4. 勾選「以最高權限執行」並取消「只在使用者登入時執行」

或使用 PowerShell（以系統管理員執行）：

```powershell
$action  = New-ScheduledTaskAction -Execute "C:\Tools\eyeProtect\daemon.exe" -Argument "-i 20 -w 20"
$trigger = New-ScheduledTaskTrigger -AtLogOn
Register-ScheduledTask -TaskName "eyeProtect" -Action $action -Trigger $trigger -RunLevel Highest
```

### macOS — launchd

建立 `~/Library/LaunchAgents/com.eyeprotect.daemon.plist`：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>             <string>com.eyeprotect.daemon</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/daemon</string>
    <string>-i</string><string>20</string>
    <string>-w</string><string>20</string>
  </array>
  <key>RunAtLoad</key>         <true/>
  <key>KeepAlive</key>         <true/>
</dict>
</plist>
```

```bash
launchctl load ~/Library/LaunchAgents/com.eyeprotect.daemon.plist
```

### Linux — systemd（使用者層級）

建立 `~/.config/systemd/user/eyeprotect.service`：

```ini
[Unit]
Description=eyeProtect 護眼提醒 Daemon
After=graphical-session.target

[Service]
ExecStart=/usr/local/bin/daemon -i 20 -w 20
Restart=on-failure

[Install]
WantedBy=default.target
```

```bash
systemctl --user enable --now eyeprotect.service
```

---

## 已知限制

| 項目 | 說明 |
|------|------|
| 舊顯卡相容性 | Intel HD 4000 等舊型整合顯卡在 Windows 上的 DX11/12 驅動不穩定，程式已強制使用 OpenGL 並停用外部 Shader 優化以確保相容 |
| 字體顯示 | 優先使用系統字體（Windows：微軟正黑；macOS：PingFang）；若找不到系統字體則回退至內建英文字體，中文可能無法正確顯示 |
| Linux 多螢幕 | 全螢幕僅作用於主螢幕，多螢幕場景視窗管理器行為因環境而異 |
| 無系統匣圖示 | 目前版本不提供系統匣（tray）控制介面，需透過命令列終止 daemon |

---

## 手動停止 Daemon

```bash
# Linux / macOS
pkill daemon

# Windows（命令提示字元）
taskkill /IM daemon.exe /F

# Windows（PowerShell）
Stop-Process -Name daemon
```

Daemon 收到終止訊號後會同時關閉正在顯示的護眼視窗，再乾淨退出。
