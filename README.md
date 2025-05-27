# Logs - High-Performance Go Logging Library

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

A fast, flexible, and feature-rich logging library for Go applications built on top of Go's standard `slog` package.

This library provides a unified interface for structured logging with advanced features like timer-based logging, stack logging, and context integration.

## Features

- **High Performance**: Outperforms popular logging libraries like Logrus, Zap, and standard slog
- **Structured Logging**: Built on Go's standard `slog` package for consistent structured output
- **Multiple Logger Types**:
  - **Standard Logger**: Fast, efficient logging with configurable levels
  - **Timer Logger**: Rate-limited logging that only outputs after specified intervals
  - **Stack Logger**: Accumulates field values into stacks for trace-like functionality
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
    logger.WithField("user_id", "12345").Info("User logged in")

    // Chain multiple fields
    logger.WithFields(map[string]interface{}{
        "method": "POST",
        "path":   "/api/users",
        "status": 201,
    }).Info("Request completed")
}
```

### Context Integration

```go
func handleRequest(ctx context.Context) {
    // Attach logger to context
    logger := logs.New(logs.LevelInfo).WithField("request_id", "abc123")
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

### Timer Logger

Timer logger is useful for preventing log spam by only outputting messages after a specified interval:

```go
// Only log once every 5 seconds
timerLogger := logs.NewTimerLogger(5*time.Second, logs.LevelInfo)

for i := 0; i < 1000; i++ {
    // This will only output a few times despite being called 1000 times
    timerLogger.Info("Processing item", i)
    time.Sleep(10 * time.Millisecond)
}
```

### Stack Logger

Stack logger accumulates values for specified field keys, creating a trace-like output:

```go
logger := logs.New(logs.LevelInfo)
stackLogger := logs.NewStackLogger(logger, "trace", "operation")

// Build up the stack
stackLogger = stackLogger.WithField("trace", "start")
stackLogger = stackLogger.WithField("operation", "validate")
stackLogger = stackLogger.WithField("trace", "middle")
stackLogger = stackLogger.WithField("operation", "process")

// Output will show: trace="start -> middle", operation="validate -> process"
stackLogger.Info("Operation completed")
```

## Logger Types

### Standard Logger

- Fast, efficient logging
- Built on Go's `slog` package
- Configurable output destination
- Support for all standard log levels

### Timer Logger

- Prevents log flooding by rate-limiting output
- Configurable time intervals
- Thread-safe implementation
- Useful for high-frequency operations

### Stack Logger

- Accumulates field values into stacks
- Great for tracing execution paths
- Configurable stack field keys
- Maintains call hierarchy

## Performance Benchmarks

Performance comparison with other popular Go logging libraries:

```bash
go test -bench=. -run=none -benchmem -v --count=1 ./test/
goos: darwin
goarch: arm64
pkg: github.com/yanun0323/logs/test
cpu: Apple M2
BenchmarkLogsBasic
BenchmarkLogsBasic-8                      582403              2038 ns/op             216 B/op          7 allocs/op
BenchmarkLogsTicker
BenchmarkLogsTicker-8                     666500              2339 ns/op             264 B/op          7 allocs/op
BenchmarkLogsTrace
BenchmarkLogsTrace-8                      464137              2360 ns/op            1096 B/op         16 allocs/op
BenchmarkSlogWithTextHandler
BenchmarkSlogWithTextHandler-8            551774              2125 ns/op             240 B/op          6 allocs/op
BenchmarkSlogWithJSONHandler
BenchmarkSlogWithJSONHandler-8            544542              2079 ns/op             240 B/op          6 allocs/op
BenchmarkSlogLogsHandler
BenchmarkSlogLogsHandler-8                585030              1983 ns/op             184 B/op          5 allocs/op
BenchmarkZap
BenchmarkZap-8                            525129              2133 ns/op            1386 B/op          8 allocs/op
BenchmarkLogrus
BenchmarkLogrus-8                         442356              2584 ns/op            1181 B/op         18 allocs/op
PASS
ok      github.com/yanun0323/logs/test  9.759s
```

**Key Performance Highlights:**

- **Fastest**: ~1.85x faster than Logrus
- **Memory Efficient**: Lower memory allocations compared to Zap and Logrus
- **Competitive**: Similar performance to standard slog while providing additional features

## API Reference

### Logger Interface

```go
type Logger interface {
    // Basic logging methods
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

    // Generic logging
    Log(level Level, args ...interface{})
    Logf(level Level, format string, args ...interface{})

    // Field management
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    WithContext(ctx context.Context) Logger

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

### Configuration

```go
// Set default logger
logger := logs.New(logs.LevelInfo)
logs.SetDefault(logger)

// Set custom time format
logs.SetDefaultTimeFormat("2006-01-02 15:04:05")
```

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
