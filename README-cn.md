# Logs - 高性能 Go 日志记录库

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

一个基于 Go 标准`slog`包构建的快速、灵活且功能丰富的 Go 应用程序日志记录库。

该库为结构化日志记录提供统一接口，具有基于计时器的日志记录、堆栈日志记录和上下文集成等高级功能。

## 特色功能

- **高性能**：性能优于 Logrus、Zap 和标准 slog 等流行日志记录库
- **结构化日志记录**：基于 Go 标准`slog`包，提供一致的结构化输出
- **多种日志记录器类型**：
  - **标准日志记录器**：快速、高效的日志记录，可配置记录级别
  - **计时器日志记录器**：速率限制的日志记录，仅在指定间隔后输出
  - **堆栈日志记录器**：将字段值累积到堆栈中，提供跟踪功能
- **上下文集成**：对 Go 上下文的一流支持
- **字段链式调用**：流畅的 API 用于添加结构化字段
- **可配置输出**：支持自定义输出写入器
- **基于级别的日志记录**：调试、信息、警告、错误和致命错误级别
- **线程安全**：安全支持并发使用

## 安装

```bash
go get github.com/yanun0323/logs
```

## 快速开始

### 基本使用

```go
package main

import (
    "github.com/yanun0323/logs"
)

func main() {
    // 使用默认日志记录器
    logs.Info("应用程序已启动")
    logs.Error("发生错误")

    // 创建自定义日志记录器
    logger := logs.New(logs.LevelDebug)
    logger.WithField("user_id", "12345").Info("用户已登录")

    // 链式添加多个字段
    logger.WithFields(map[string]interface{}{
        "method": "POST",
        "path":   "/api/users",
        "status": 201,
    }).Info("请求已完成")
}
```

### 上下文集成

```go
func handleRequest(ctx context.Context) {
    // 将日志记录器附加到上下文
    logger := logs.New(logs.LevelInfo).WithField("request_id", "abc123")
    ctx = logger.Attach(ctx)

    // 在后续的调用链中
    processData(ctx)
}

func processData(ctx context.Context) {
    // 从上下文中获取日志记录器
    logger := logs.Get(ctx)
    logger.Info("正在处理数据")
}
```

### 计时器日志记录器

计时器日志记录器对于防止日志洪水很有用，只有在指定间隔后才输出消息：

```go
// 每5秒只记录一次
timerLogger := logs.NewTimerLogger(5*time.Second, logs.LevelInfo)

for i := 0; i < 1000; i++ {
    // 尽管被调用1000次，但只会输出几次
    timerLogger.Info("正在处理项目", i)
    time.Sleep(10 * time.Millisecond)
}
```

### 堆栈日志记录器

堆栈日志记录器为指定的字段键累积值，创建类似跟踪的输出：

```go
logger := logs.New(logs.LevelInfo)
stackLogger := logs.NewStackLogger(logger, "trace", "operation")

// 构建堆栈
stackLogger = stackLogger.WithField("trace", "start")
stackLogger = stackLogger.WithField("operation", "validate")
stackLogger = stackLogger.WithField("trace", "middle")
stackLogger = stackLogger.WithField("operation", "process")

// 输出将显示：trace="start -> middle", operation="validate -> process"
stackLogger.Info("操作已完成")
```

## 日志记录器类型

### 标准日志记录器

- 快速、高效的日志记录
- 基于 Go 的`slog`包
- 可配置输出目标
- 支持所有标准日志级别

### 计时器日志记录器

- 通过速率限制防止日志洪水
- 可配置时间间隔
- 线程安全实现
- 适用于高频率操作

### 堆栈日志记录器

- 将字段值累积到堆栈中
- 非常适合跟踪执行路径
- 可配置堆栈字段键
- 维护调用层次结构

## 性能基准测试

与其他流行 Go 日志记录库的性能比较：

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

**关键性能亮点：**

- **最快速**：比 Logrus 快约 1.85 倍
- **内存高效**：相比 Zap 和 Logrus 的内存分配更少
- **竞争力强**：在提供额外功能的同时，性能与标准 slog 相似

## API 参考

### 日志记录器接口

```go
type Logger interface {
    // 基本日志记录方法
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

    // 通用日志记录
    Log(level Level, args ...interface{})
    Logf(level Level, format string, args ...interface{})

    // 字段管理
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    WithContext(ctx context.Context) Logger

    // 工具方法
    Copy() Logger
    Attach(ctx context.Context) context.Context
}
```

### 日志级别

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
// 设置默认日志记录器
logger := logs.New(logs.LevelInfo)
logs.SetDefault(logger)

// 设置自定义时间格式
logs.SetDefaultTimeFormat("2006-01-02 15:04:05")
```

## 依赖项

该库的依赖项很少，仅使用 Go 标准库：

- **Go 1.21 或更高版本** - 所需的 Go 版本
- **仅使用标准库** - 核心功能无需外部依赖
- **可选的测试依赖项**：
  - `github.com/sirupsen/logrus` - 用于性能基准测试
  - `go.uber.org/zap` - 用于性能基准测试

该库设计为轻量且无依赖，同时提供强大的日志记录功能。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件。
