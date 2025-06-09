# Logs - 高效能 Go 日誌記錄庫

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

一個基於 Go 標準`slog`套件建構的快速、靈活且功能豐富的 Go 應用程式日誌記錄庫。

本庫提供了結構化日誌記錄的統一介面，具有計時器式日誌記錄、堆疊日誌記錄和上下文整合等進階功能。

## 外觀

### Console 格式

![Console Format](https://github.com/yanun0323/assets/blob/master/logs.appearance.png?raw=true)

### 文字格式

```log
time=2025-05-28T04:29:22.422+08:00 level=DEBUG msg="debug message"
time=2025-05-28T04:29:22.422+08:00 level=INFO msg="info message with fields" fields=val error=<nil> context=context.TODO func=testFunc
time=2025-05-28T04:29:22.422+08:00 level=WARN msg="warn message with func trace" func="testFunc -> testFunc2 -> testFunc3"
time=2025-05-28T04:29:22.422+08:00 level=ERROR msg="error message"
time=2025-05-28T04:29:22.422+08:00 level=ERROR+4 msg="fatal message"
```

### JSON 格式

```log
{"time":"2025-05-28T04:24:56.279024+08:00","level":"DEBUG","msg":"debug message"}
{"time":"2025-05-28T04:24:56.279113+08:00","level":"INFO","msg":"info message with fields","fields":"val","error":null,"context":{},"func":"testFunc"}
{"time":"2025-05-28T04:24:56.279127+08:00","level":"WARN","msg":"warn message with func trace","func":"testFunc -> testFunc2 -> testFunc3"}
{"time":"2025-05-28T04:24:56.279137+08:00","level":"ERROR","msg":"error message"}
{"time":"2025-05-28T04:24:56.279139+08:00","level":"ERROR+4","msg":"fatal message"}
```

## 特色功能

- **高效能**：效能優於 Logrus、Zap 和標準 slog 等熱門日誌記錄庫
- **結構化日誌記錄**：基於 Go 標準`slog`套件，提供一致的結構化輸出
- **多種日誌記錄器類型**：
  - **標準日誌記錄器**：快速、高效的日誌記錄，可配置記錄等級
  - **TickerLogger**：速率限制的日誌記錄，僅在指定間隔後輸出
- **上下文整合**：對 Go 上下文的一流支援
- **欄位鏈式呼叫**：流暢的 API 用於新增結構化欄位
- **可配置輸出**：支援自訂輸出寫入器
- **基於等級的日誌記錄**：除錯、資訊、警告、錯誤和致命錯誤等級
- **執行緒安全**：安全支援並發使用

## 安裝

```bash
go get github.com/yanun0323/logs
```

## 快速開始

### 基本使用

```go
package main

import (
    "github.com/yanun0323/logs"
)

func main() {
    // 使用預設日誌記錄器
    logs.Info("應用程式已啟動")
    logs.Error("發生錯誤")

    // 建立自訂日誌記錄器
    logger := logs.New(logs.LevelDebug)
    logger.With("user_id", "12345").Info("使用者已登入")

    // 鏈式新增多個欄位
    logger.With(
        "method": "POST",
        "path": "/api/users",
        "status": 201,
    ).Info("請求已完成")
}
```

### 上下文整合

```go
func handleRequest(ctx context.Context) {
    // 將日誌記錄器附加到上下文
    logger := logs.New(logs.LevelInfo).With("request_id", "abc123")
    ctx = logger.Attach(ctx)

    // 在後續的呼叫鏈中
    processData(ctx)
}

func processData(ctx context.Context) {
    // 從上下文中取得日誌記錄器
    logger := logs.Get(ctx)
    logger.Info("正在處理資料")
}
```

### TickerLogger

TickerLogger 對於防止日誌洪水很有用，只有在指定間隔後才輸出訊息：

```go
// 每5秒只記錄一次
tickerLogger := logs.NewTickerLogger(logs.LevelInfo, 5*time.Second)

for i := 0; i < 1000; i++ {
    // 儘管被呼叫1000次，但只會輸出幾次
    tickerLogger.Info("正在處理項目", i)
    time.Sleep(10 * time.Millisecond)
}
```

## 日誌記錄器類型

### 標準日誌記錄器

- 快速、高效的日誌記錄
- 基於 Go 的`slog`套件
- 可配置輸出目標
- 支援所有標準日誌等級

### TickerLogger

- 透過速率限制防止日誌洪水
- 可配置時間間隔
- 執行緒安全實作
- 適用於高頻率操作

## 效能基準測試

與其他熱門 Go 日誌記錄庫的效能比較：

```bash
go test -bench=. -run=none -benchmem -v --count=1 -benchtime=30s ./test/
goos: darwin
goarch: arm64
pkg: github.com/yanun0323/logs/test
cpu: Apple M2
BenchmarkLogsBasic-8                     8626179              4827 ns/op             272 B/op          9 allocs/op
BenchmarkLogsTicker-8                  169353175             212.1 ns/op             272 B/op          8 allocs/op

BenchmarkSlogWithTextHandler-8           8255836              5003 ns/op             240 B/op          6 allocs/op
BenchmarkSlogWithJSONHandler-8           7863826              5220 ns/op             240 B/op          6 allocs/op
BenchmarkSlogLogsHandler-8               9155694              4654 ns/op             224 B/op          6 allocs/op
BenchmarkZap-8                           7550286              5410 ns/op            1410 B/op         10 allocs/op
BenchmarkLogrus-8                        6633025              6055 ns/op            1593 B/op         31 allocs/op
PASS
ok      github.com/yanun0323/logs/test  323.964s
```

**關鍵效能亮點：**

- **效能優於競爭對手**：基本日誌記錄器比 Zap 快約 1.25 倍，比 Logrus 快約 1.26 倍
- **超高速 TickerLogger**：比基本日誌記錄快約 23 倍（212.1 ns/op vs 4827 ns/op）
- **記憶體高效**：相比 Zap（272 B vs 1410 B）和 Logrus（272 B vs 1593 B）的記憶體分配顯著更少
- **最快處理器**：自訂 slog 處理器達到最佳效能，為 4654 ns/op
- **與 slog 競爭力強**：在提供額外功能的同時，效能與標準 slog 相當

## API 參考

### 日誌記錄器介面

```go
type Logger interface {
    // 基本日誌記錄方法
    Debug(args ...any)
    Debugf(format string, args ...any)
    Info(args ...any)
    Infof(format string, args ...any)
    Warn(args ...any)
    Warnf(format string, args ...any)
    Error(args ...any)
    Errorf(format string, args ...any)
    Fatal(args ...any)
    Fatalf(format string, args ...any)

    // 欄位管理
    With(args ...any) Logger
    WithErr(err error) Logger
    WithCtx(ctx context.Context) Logger
    WithFunc(function string) Logger

    // 工具方法
    Copy() Logger
    Attach(ctx context.Context) context.Context
}
```

### 日誌等級

```go
const (
    LevelDebug Level = -4
    LevelInfo  Level = 0
    LevelWarn  Level = 4
    LevelError Level = 8
    LevelFatal Level = 12
)
```

### 輸出格式

本庫支援多種輸出格式以適應不同的使用情境：

```go
type Format int8

const (
    // FormatConsole 輸出彩色、人類可讀的控制台格式。
    // 這是預設格式，在開發時具有增強的可讀性。
    FormatConsole Format = iota + 1

    // FormatText 輸出標準 slog 文字格式。
    // 格式：以空格分隔的 key=value 對。
    FormatText

    // FormatJSON 輸出 JSON 格式。
    // 每個日誌條目都是一行中的單個 JSON 物件。
    FormatJSON
)
```

### 日誌記錄器選項

使用 `Option` 結構自訂日誌記錄器行為：

```go
type Option struct {
    // Format 指定日誌輸出格式。
    // 可用格式：FormatConsole（預設）、FormatText、FormatJSON
    Format Format

    // Output 指定日誌輸出的目標寫入器。
    // 如果未指定，預設為 os.Stdout。
    Output io.Writer
}
```

**使用範例：**

```go
// 建立 JSON 格式的日誌記錄器
logger := logs.New(logs.LevelInfo, &logs.Option{
    Format: logs.FormatJSON,
    Output: os.Stdout,
})

// 建立文字格式並寫入檔案的日誌記錄器
file, _ := os.Create("app.log")
logger := logs.New(logs.LevelDebug, &logs.Option{
    Format: logs.FormatText,
    Output: file,
})

// 建立控制台格式的日誌記錄器（預設）
logger := logs.New(logs.LevelInfo, &logs.Option{
    Format: logs.FormatConsole,
    Output: os.Stdout,
})

// 或簡單使用預設值
logger := logs.New(logs.LevelInfo)
```

### 配置

```go
// 設定預設日誌記錄器
logger := logs.New(logs.LevelInfo)
logs.SetDefault(logger)

// 設定自訂時間格式
logs.SetDefaultTimeFormat("2006-01-02 15:04:05")
```

## 依賴項目

本庫的依賴項目很少，僅使用 Go 標準庫：

- **Go 1.21 或更高版本** - 所需的 Go 版本
- **僅使用標準庫** - 核心功能無需外部依賴
- **可選的測試依賴項目**：
  - `github.com/sirupsen/logrus` - 用於效能基準測試
  - `go.uber.org/zap` - 用於效能基準測試

本庫設計為輕量且無依賴，同時提供強大的日誌記錄功能。

## 貢獻

歡迎貢獻！請隨時提交 Pull Request。

## 授權條款

本專案採用 MIT 授權條款 - 詳情請參閱 LICENSE 檔案。
