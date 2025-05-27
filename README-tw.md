# Logs - 高效能 Go 日誌記錄庫

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

一個基於 Go 標準`slog`套件建構的快速、靈活且功能豐富的 Go 應用程式日誌記錄庫。

本庫提供了結構化日誌記錄的統一介面，具有計時器式日誌記錄、堆疊日誌記錄和上下文整合等進階功能。

## 特色功能

- **高效能**：效能優於 Logrus、Zap 和標準 slog 等熱門日誌記錄庫
- **結構化日誌記錄**：基於 Go 標準`slog`套件，提供一致的結構化輸出
- **多種日誌記錄器類型**：
  - **標準日誌記錄器**：快速、高效的日誌記錄，可配置記錄等級
  - **計時器日誌記錄器**：速率限制的日誌記錄，僅在指定間隔後輸出
  - **堆疊日誌記錄器**：將欄位值累積到堆疊中，提供追蹤功能
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
    logger.WithField("user_id", "12345").Info("使用者已登入")

    // 鏈式新增多個欄位
    logger.WithFields(map[string]interface{}{
        "method": "POST",
        "path":   "/api/users",
        "status": 201,
    }).Info("請求已完成")
}
```

### 上下文整合

```go
func handleRequest(ctx context.Context) {
    // 將日誌記錄器附加到上下文
    logger := logs.New(logs.LevelInfo).WithField("request_id", "abc123")
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

### 計時器日誌記錄器

計時器日誌記錄器對於防止日誌洪水很有用，只有在指定間隔後才輸出訊息：

```go
// 每5秒只記錄一次
timerLogger := logs.NewTimerLogger(5*time.Second, logs.LevelInfo)

for i := 0; i < 1000; i++ {
    // 儘管被呼叫1000次，但只會輸出幾次
    timerLogger.Info("正在處理項目", i)
    time.Sleep(10 * time.Millisecond)
}
```

### 堆疊日誌記錄器

堆疊日誌記錄器為指定的欄位鍵累積值，建立類似追蹤的輸出：

```go
logger := logs.New(logs.LevelInfo)
stackLogger := logs.NewStackLogger(logger, "trace", "operation")

// 建立堆疊
stackLogger = stackLogger.WithField("trace", "start")
stackLogger = stackLogger.WithField("operation", "validate")
stackLogger = stackLogger.WithField("trace", "middle")
stackLogger = stackLogger.WithField("operation", "process")

// 輸出將顯示：trace="start -> middle", operation="validate -> process"
stackLogger.Info("操作已完成")
```

## 日誌記錄器類型

### 標準日誌記錄器

- 快速、高效的日誌記錄
- 基於 Go 的`slog`套件
- 可配置輸出目標
- 支援所有標準日誌等級

### 計時器日誌記錄器

- 透過速率限制防止日誌洪水
- 可配置時間間隔
- 執行緒安全實作
- 適用於高頻率操作

### 堆疊日誌記錄器

- 將欄位值累積到堆疊中
- 非常適合追蹤執行路徑
- 可配置堆疊欄位鍵
- 維護呼叫階層

## 效能基準測試

與其他熱門 Go 日誌記錄庫的效能比較：

```bash
$ go test -bench=. -run=none -benchmem --count=1 ./test
goos: darwin
goarch: arm64
pkg: github.com/yanun0323/logs/test
cpu: Apple M2
BenchmarkLogs-8                           469046              2533 ns/op             600 B/op         20 allocs/op
BenchmarkLogrus-8                         259897              4149 ns/op            1181 B/op         18 allocs/op
BenchmarkSlogWithLogsHandler-8             365598              3958 ns/op             584 B/op         18 allocs/op
BenchmarkSlog-8                           344457              3430 ns/op             240 B/op          6 allocs/op
BenchmarkZap-8                            327561              3551 ns/op            1388 B/op          8 allocs/op
PASS
ok      github.com/yanun0323/logs/test  7.305s
```

**關鍵效能亮點：**

- **最快速**：比 Logrus 快約 1.85 倍
- **記憶體高效**：相比 Zap 和 Logrus 的記憶體分配更少
- **競爭力強**：在提供額外功能的同時，效能與標準 slog 相似

## API 參考

### 日誌記錄器介面

```go
type Logger interface {
    // 基本日誌記錄方法
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    Warn(args ...interface{})
    Warnf(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Fatal(args ...interface{})
    Fatalf(format string, args ...interface{})

    // 通用日誌記錄
    Log(level Level, args ...interface{})
    Logf(level Level, format string, args ...interface{})

    // 欄位管理
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    WithContext(ctx context.Context) Logger

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
