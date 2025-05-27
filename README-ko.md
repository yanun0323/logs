# Logs - 고성능 Go 로깅 라이브러리

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

Go의 표준 `slog` 패키지를 기반으로 구축된 빠르고 유연하며 기능이 풍부한 Go 애플리케이션용 로깅 라이브러리입니다.

이 라이브러리는 타이머 기반 로깅, 스택 로깅, 컨텍스트 통합과 같은 고급 기능을 갖춘 구조화된 로깅을 위한 통합 인터페이스를 제공합니다.

## 특징

- **고성능**: Logrus, Zap, 표준 slog와 같은 인기 있는 로깅 라이브러리보다 뛰어난 성능
- **구조화된 로깅**: 일관된 구조화된 출력을 위한 Go의 표준 `slog` 패키지 기반
- **다양한 로거 타입**:
  - **표준 로거**: 구성 가능한 레벨을 가진 빠르고 효율적인 로깅
  - **타이머 로거**: 지정된 간격 후에만 출력하는 속도 제한 로깅
  - **스택 로거**: 필드 값을 스택에 누적하여 추적 기능 제공
- **컨텍스트 통합**: Go 컨텍스트에 대한 일급 지원
- **필드 체이닝**: 구조화된 필드 추가를 위한 플루언트 API
- **구성 가능한 출력**: 사용자 정의 출력 라이터 지원
- **레벨 기반 로깅**: 디버그, 정보, 경고, 오류, 치명적 오류 레벨
- **스레드 안전**: 동시 사용에 안전

## 설치

```bash
go get github.com/yanun0323/logs
```

## 빠른 시작

### 기본 사용법

```go
package main

import (
    "github.com/yanun0323/logs"
)

func main() {
    // 기본 로거 사용
    logs.Info("애플리케이션이 시작되었습니다")
    logs.Error("문제가 발생했습니다")

    // 사용자 정의 로거 생성
    logger := logs.New(logs.LevelDebug)
    logger.WithField("user_id", "12345").Info("사용자가 로그인했습니다")

    // 여러 필드 체이닝
    logger.WithFields(map[string]interface{}{
        "method": "POST",
        "path":   "/api/users",
        "status": 201,
    }).Info("요청이 완료되었습니다")
}
```

### 컨텍스트 통합

```go
func handleRequest(ctx context.Context) {
    // 로거를 컨텍스트에 연결
    logger := logs.New(logs.LevelInfo).WithField("request_id", "abc123")
    ctx = logger.Attach(ctx)

    // 후속 호출 체인에서
    processData(ctx)
}

func processData(ctx context.Context) {
    // 컨텍스트에서 로거 가져오기
    logger := logs.Get(ctx)
    logger.Info("데이터 처리 중")
}
```

### 타이머 로거

타이머 로거는 지정된 간격 후에만 메시지를 출력하여 로그 플러딩을 방지하는 데 유용합니다:

```go
// 5초마다 한 번만 로그 기록
timerLogger := logs.NewTimerLogger(5*time.Second, logs.LevelInfo)

for i := 0; i < 1000; i++ {
    // 1000번 호출되어도 몇 번만 출력됩니다
    timerLogger.Info("항목 처리 중", i)
    time.Sleep(10 * time.Millisecond)
}
```

### 스택 로거

스택 로거는 지정된 필드 키에 대한 값을 누적하여 추적과 같은 출력을 생성합니다:

```go
logger := logs.New(logs.LevelInfo)
stackLogger := logs.NewStackLogger(logger, "trace", "operation")

// 스택 구축
stackLogger = stackLogger.WithField("trace", "start")
stackLogger = stackLogger.WithField("operation", "validate")
stackLogger = stackLogger.WithField("trace", "middle")
stackLogger = stackLogger.WithField("operation", "process")

// 출력 표시: trace="start -> middle", operation="validate -> process"
stackLogger.Info("작업이 완료되었습니다")
```

## 로거 타입

### 표준 로거

- 빠르고 효율적인 로깅
- Go의 `slog` 패키지 기반
- 구성 가능한 출력 대상
- 모든 표준 로그 레벨 지원

### 타이머 로거

- 속도 제한을 통한 로그 플러딩 방지
- 구성 가능한 시간 간격
- 스레드 안전 구현
- 고빈도 작업에 유용

### 스택 로거

- 필드 값을 스택에 누적
- 실행 경로 추적에 뛰어남
- 구성 가능한 스택 필드 키
- 호출 계층 구조 유지

## 성능 벤치마크

다른 인기 있는 Go 로깅 라이브러리와의 성능 비교:

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

**주요 성능 하이라이트:**

- **가장 빠름**: Logrus보다 약 1.85배 빠름
- **메모리 효율적**: Zap과 Logrus에 비해 적은 메모리 할당
- **경쟁력**: 추가 기능을 제공하면서 표준 slog와 유사한 성능

## API 참조

### 로거 인터페이스

```go
type Logger interface {
    // 기본 로깅 메서드
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

    // 범용 로깅
    Log(level Level, args ...interface{})
    Logf(level Level, format string, args ...interface{})

    // 필드 관리
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    WithContext(ctx context.Context) Logger

    // 유틸리티 메서드
    Copy() Logger
    Attach(ctx context.Context) context.Context
}
```

### 로그 레벨

```go
const (
    LevelDebug Level = -4
    LevelInfo  Level = 0
    LevelWarn  Level = 4
    LevelError Level = 8
    LevelFatal Level = 12
)
```

### 구성

```go
// 기본 로거 설정
logger := logs.New(logs.LevelInfo)
logs.SetDefault(logger)

// 사용자 정의 시간 형식 설정
logs.SetDefaultTimeFormat("2006-01-02 15:04:05")
```

## 의존성

이 라이브러리는 최소한의 의존성을 가지며 Go 표준 라이브러리만 사용합니다:

- **Go 1.21 이상** - 필요한 Go 버전
- **표준 라이브러리만 사용** - 핵심 기능에 외부 의존성이 없음
- **선택적 테스트 의존성**:
  - `github.com/sirupsen/logrus` - 성능 벤치마킹용
  - `go.uber.org/zap` - 성능 벤치마킹용

이 라이브러리는 강력한 로깅 기능을 제공하면서도 가볍고 의존성이 없는 설계입니다.

## 기여

기여를 환영합니다! 언제든지 Pull Request를 제출해 주세요.

## 라이선스

이 프로젝트는 MIT 라이선스 하에 라이선스가 부여됩니다 - 자세한 내용은 LICENSE 파일을 참조하십시오.
