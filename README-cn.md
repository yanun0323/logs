# Logs - 高性能 Go 日志记录库

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

一个基于 Go 标准`slog`包构建的快速、灵活且功能丰富的 Go 应用程序日志记录库。

该库为结构化日志记录提供统一接口，具有基于计时器的日志记录、堆栈日志记录和上下文集成等高级功能。

## 外观

### Console 格式

![Console Format](https://github.com/yanun0323/assets/blob/master/logs.appearance.png?raw=true)

### 文本格式

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

- **高性能**：性能优于 Logrus、Zap 和标准 slog 等流行日志记录库
- **结构化日志记录**：基于 Go 标准`slog`包，提供一致的结构化输出
- **多种日志记录器类型**：
  - **标准日志记录器**：快速、高效的日志记录，可配置记录级别
  - **TickerLogger**：速率限制的日志记录，仅在指定间隔后输出
  - **TraceLogger**：将字段值累积到堆栈中，提供跟踪功能
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
    logger.WithFields(map[string]any{
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

### TickerLogger

TickerLogger 对于防止日志洪水很有用，只有在指定间隔后才输出消息：

```go
// 每5秒只记录一次
tickerLogger := logs.NewTickerLogger(logs.LevelInfo, 5*time.Second)

for i := 0; i < 1000; i++ {
    // 尽管被调用1000次，但只会输出几次
    tickerLogger.Info("正在处理项目", i)
    time.Sleep(10 * time.Millisecond)
}
```

### TraceLogger

TraceLogger 为指定的字段键累积值，创建类似跟踪的输出：

```go
// 构建堆栈
traceLogger := logs.NewTraceLogger(logs.LevelInfo, "trace")
traceLogger = traceLogger.WithField("trace", "start")
traceLogger = traceLogger.WithField("trace", "middle")
traceLogger = traceLogger.WithField("trace", "end")

// 输出将显示：trace="start -> middle -> end"
traceLogger.Info("操作已完成")
```

## 日志记录器类型

### 标准日志记录器

- 快速、高效的日志记录
- 基于 Go 的`slog`包
- 可配置输出目标
- 支持所有标准日志级别

### TickerLogger

- 通过速率限制防止日志洪水
- 可配置时间间隔
- 线程安全实现
- 适用于高频率操作

### TraceLogger

- 将字段值累积到堆栈中
- 非常适合跟踪执行路径
- 可配置堆栈字段键
- 维护调用层次结构

## 性能基准测试

与其他流行 Go 日志记录库的性能比较：

```bash
go test -bench=. -run=none -benchmem -v --count=1 -benchtime=30s ./test/
goos: darwin
goarch: arm64
pkg: github.com/yanun0323/logs/test
cpu: Apple M2
BenchmarkLogsBasic-8                     8626179              4827 ns/op             272 B/op          9 allocs/op
BenchmarkLogsTicker-8                  169353175             212.1 ns/op             272 B/op          8 allocs/op
BenchmarkLogsTrace-8                     7778780              5162 ns/op            1152 B/op         18 allocs/op
BenchmarkSlogWithTextHandler-8           8255836              5003 ns/op             240 B/op          6 allocs/op
BenchmarkSlogWithJSONHandler-8           7863826              5220 ns/op             240 B/op          6 allocs/op
BenchmarkSlogLogsHandler-8               9155694              4654 ns/op             224 B/op          6 allocs/op
BenchmarkZap-8                           7550286              5410 ns/op            1410 B/op         10 allocs/op
BenchmarkLogrus-8                        6633025              6055 ns/op            1593 B/op         31 allocs/op
PASS
ok      github.com/yanun0323/logs/test  323.964s
```

**关键性能亮点：**

- **性能优于竞争对手**：基本日志记录器比 Zap 快约 1.25 倍，比 Logrus 快约 1.26 倍
- **超高速 TickerLogger**：比基本日志记录快约 23 倍（212.1 ns/op vs 4827 ns/op）
- **内存高效**：相比 Zap（272 B vs 1410 B）和 Logrus（272 B vs 1593 B）的内存分配显著更少
- **最快处理器**：自定义 slog 处理器达到最佳性能，为 4654 ns/op
- **与 slog 竞争力强**：在提供额外功能的同时，性能与标准 slog 相当

## API 参考

### 日志记录器接口

```go
type Logger interface {
    // 基本日志记录方法
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

    // 字段管理
    WithField(key string, value any) Logger
    WithFields(args ...any) Logger
    WithError(err error) Logger
    WithContext(ctx context.Context) Logger
    WithFunc(function string) Logger

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

### 输出格式

该库支持多种输出格式以适应不同的使用场景：

```go
type Format int8

const (
    // FormatConsole 输出彩色、人类可读的控制台格式。
    // 这是默认格式，在开发时具有增强的可读性。
    FormatConsole Format = iota + 1

    // FormatText 输出标准 slog 文本格式。
    // 格式：以空格分隔的 key=value 对。
    FormatText

    // FormatJSON 输出 JSON 格式。
    // 每个日志条目都是一行中的单个 JSON 对象。
    FormatJSON
)
```

### 日志记录器选项

使用 `Option` 结构自定义日志记录器行为：

```go
type Option struct {
    // Format 指定日志输出格式。
    // 可用格式：FormatConsole（默认）、FormatText、FormatJSON
    Format Format

    // Output 指定日志输出的目标写入器。
    // 如果未指定，默认为 os.Stdout。
    Output io.Writer
}
```

**使用示例：**

```go
// 创建 JSON 格式的日志记录器
logger := logs.New(logs.LevelInfo, &logs.Option{
    Format: logs.FormatJSON,
    Output: os.Stdout,
})

// 创建文本格式并写入文件的日志记录器
file, _ := os.Create("app.log")
logger := logs.New(logs.LevelDebug, &logs.Option{
    Format: logs.FormatText,
    Output: file,
})

// 创建控制台格式的日志记录器（默认）
logger := logs.New(logs.LevelInfo, &logs.Option{
    Format: logs.FormatConsole,
    Output: os.Stdout,
})

// 或简单使用默认值
logger := logs.New(logs.LevelInfo)
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
