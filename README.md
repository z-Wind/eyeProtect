# eyeProtect

## 編譯方法 (Install)

1. 編譯顯示端:
   go build -o eyeProtect.exe ./cmd/ui

2. 編譯背景主控端:
   go build -o daemon.exe ./cmd/daemon

## 參數說明 (Flags)

  -t, --top-enable                           開啟此選項後，護眼視窗將會覆蓋在所有視窗之上 [default: true]

  -w, --wait-seconds <WAIT_SECONDS>          設定護眼視窗出現的時間長度（單位：秒） [default: 20]

  -i, --interval-minutes <INTERVAL_MINUTES>  設定每隔多久提醒一次休息（單位：分鐘） [default: 10]

  -r, --remind <REMIND>                      設定自定義的提醒文字（例如：喝口水、站起來動一動）


## 執行範例 (Run)

### Windows (背景執行)
daemon.exe -i 15 -w 30 -r "該看窗外囉" -t

### Linux
./daemon -i 10 -w 20 -r "站起來動一動"
