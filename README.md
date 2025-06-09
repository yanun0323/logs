# Logs - High-Performance Go Logging Library

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

A fast, flexible, and feature-rich logging library for Go applications built on top of Go's standard `slog` package.

This library provides a unified interface for structured logging with advanced features like timer-based logging, stack logging, and context integration.

## Appearance

#### Console Format

![Console Format](https://github.com/yanun0323/assets/blob/master/logs.appearance.png?raw=true)

#### Text Format

```log
time=2025-05-28T04:29:22.422+08:00 level=DEBUG msg="debug message"
time=2025-05-28T04:29:22.422+08:00 level=INFO msg="info message with fields" fields=val error=<nil> context=context.TODO func=testFunc
time=2025-05-28T04:29:22.422+08:00 level=WARN msg="warn message with func trace" func="testFunc -> testFunc2 -> testFunc3"
time=2025-05-28T04:29:22.422+08:00 level=ERROR msg="error message"
time=2025-05-28T04:29:22.422+08:00 level=ERROR+4 msg="fatal message"
```

#### JSON Format

```log
{"time":"2025-05-28T04:24:56.279024+08:00","level":"DEBUG","msg":"debug message"}
{"time":"2025-05-28T04:24:56.279113+08:00","level":"INFO","msg":"info message with fields","fields":"val","error":null,"context":{},"func":"testFunc"}
{"time":"2025-05-28T04:24:56.279127+08:00","level":"WARN","msg":"warn message with func trace","func":"testFunc -> testFunc2 -> testFunc3"}
{"time":"2025-05-28T04:24:56.279137+08:00","level":"ERROR","msg":"error message"}
{"time":"2025-05-28T04:24:56.279139+08:00","level":"ERROR+4","msg":"fatal message"}
```

## Features

- **High Performance**: Outperforms popular logging libraries like Logrus, Zap, and standard slog
- **Structured Logging**: Built on Go's standard `slog` package for consistent structured output
- **Multiple Logger Types**:
  - **Standard Logger**: Fast, efficient logging with configurable levels
  - **Ticker Logger**: Rate-limited logging that only outputs after specified intervals
- **Context Integration**: First-class support for Go contexts
- **Field Chaining**: Fluent API for adding structured fields
- **Configurable Output**: Support for custom output writers
- **Level-based Logging**: Debug, Info, Warn, Error, and Fatal levels
- **Thread-safe**: Safe for concurrent use

## Installation

```bash
go get github.com/yanun0323/logs
```

## Quick Start

### Basic Usage

```go
package main

import (
    "github.com/yanun0323/logs"
)

func main() {
    // Use default logger
    logs.Info("Application started")
    logs.Error("Something went wrong")

    // Create custom logger
    logger := logs.New(logs.LevelDebug)
    logger.With("user_id", "12345").Info("User logged in")

    // Chain multiple fields
    logger.With(
        "method", "POST",
        "path", "/api/users",
        "status", 201,
    ).Info("Request completed")
}
```

### Context Integration

```go
func handleRequest(ctx context.Context) {
    // Attach logger to context
    logger := logs.New(logs.LevelInfo).With("request_id", "abc123")
    ctx = logger.Attach(ctx)

    // Later in the call chain
    processData(ctx)
}

func processData(ctx context.Context) {
    // Get logger from context
    logger := logs.Get(ctx)
    logger.Info("Processing data")
}
```

### Ticker Logger

Ticker logger is useful for preventing log spam by only outputting messages after a specified interval:

```go
// Only log once every 5 seconds
tickerLogger := logs.NewTickerLogger(logs.LevelInfo, 5*time.Second)

for i := 0; i < 1000; i++ {
    // This will only output a few times despite being called 1000 times
    tickerLogger.Info("Processing item", i)
    time.Sleep(10 * time.Millisecond)
}
```

## Logger Types

### Standard Logger

- Fast, efficient logging
- Built on Go's `slog` package
- Configurable output destination
- Support for all standard log levels

### Ticker Logger

- Prevents log flooding by rate-limiting output
- Configurable time intervals
- Thread-safe implementation
- Useful for high-frequency operations

## Performance Benchmarks

Performance comparison with other popular Go logging libraries:

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

**Key Performance Highlights:**

- **Outperforms Competitors**: Basic logger is ~1.25x faster than Zap and ~1.26x faster than Logrus
- **Ultra-fast Ticker Logger**: ~23x faster than basic logging (212.1 ns/op vs 4827 ns/op)
- **Memory Efficient**: Significantly lower memory allocations compared to Zap (272 B vs 1410 B) and Logrus (272 B vs 1593 B)
- **Fastest Handler**: Custom slog handler achieves best performance at 4654 ns/op
- **Competitive with slog**: Performance comparable to standard slog while providing additional features

## API Reference

### Logger Interface

```go
type Logger interface {
    // Basic logging methods
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

    // Field management
    With(args ...any) Logger
    WithErr(err error) Logger
    WithCtx(ctx context.Context) Logger
    WithFunc(function string) Logger

    // Utility methods
    Copy() Logger
    Attach(ctx context.Context) context.Context
}
```

### Log Levels

```go
const (
    LevelDebug Level = -4
    LevelInfo  Level = 0
    LevelWarn  Level = 4
    LevelError Level = 8
    LevelFatal Level = 12
)
```

### Output Formats

The library supports multiple output formats for different use cases:

```go
type Format int8

const (
    // FormatConsole outputs logs in a colored, human-readable console format.
    // This is the default format with enhanced readability for development.
    FormatConsole Format = iota + 1

    // FormatText outputs logs in standard slog text format.
    // Format: key=value pairs separated by spaces.
    FormatText

    // FormatJSON outputs logs in JSON format.
    // Each log entry is a single JSON object on one line.
    FormatJSON
)
```

### Logger Options

Customize logger behavior using the `Option` struct:

```go
type Option struct {
    // Format specifies the log output format.
    // Available formats: FormatConsole (default), FormatText, FormatJSON
    Format Format

    // Output specifies the destination writer for log output.
    // Defaults to os.Stdout if not specified.
    Output io.Writer
}
```

**Usage Examples:**

```go
// Create logger with JSON format
logger := logs.New(logs.LevelInfo, &logs.Option{
    Format: logs.FormatJSON,
    Output: os.Stdout,
})

// Create logger with text format writing to file
file, _ := os.Create("app.log")
logger := logs.New(logs.LevelDebug, &logs.Option{
    Format: logs.FormatText,
    Output: file,
})

// Create logger with console format (default)
logger := logs.New(logs.LevelInfo, &logs.Option{
    Format: logs.FormatConsole,
    Output: os.Stdout,
})

// Or simply use defaults
logger := logs.New(logs.LevelInfo)
```

### Configuration

```go
// Set default logger
logger := logs.New(logs.LevelInfo)
logs.SetDefault(logger)

// Set custom time format
logs.SetDefaultTimeFormat("2006-01-02 15:04:05")
```

### Errors Package Integration

This package interoperates with the [github.com/yanun0323/logs](https://github.com/yanun0323/errors) package.

```go
logger := logs.Default()

err := errors.New("database connection failed").
    With("host", "localhost").
    With("port", 5432).

logger.WithError(err).Error("Operation error")
```

When using with the `errors` package, errors created by `errors` package can be directly passed to log functions and will automatically extract structured fields and stack traces.

## Dependencies

This library has minimal dependencies and uses only Go's standard library:

- **Go 1.21 or later** - Required Go version
- **Standard library only** - No external dependencies for core functionality
- **Optional testing dependencies**:
  - `github.com/sirupsen/logrus` - For performance benchmarking
  - `go.uber.org/zap` - For performance benchmarking

The library is designed to be lightweight and dependency-free while providing powerful logging capabilities.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
