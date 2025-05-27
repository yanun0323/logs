# Logs - 高性能 Go ロギングライブラリ

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)
[![한국어](https://img.shields.io/badge/한국어-클릭-yellow)](README-ko.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/yanun0323/logs.svg)](https://pkg.go.dev/github.com/yanun0323/logs)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanun0323/logs)](https://goreportcard.com/report/github.com/yanun0323/logs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)

Go の標準`slog`パッケージをベースに構築された、高速で柔軟且つ機能豊富な Go アプリケーション用ロギングライブラリです。

このライブラリは、タイマーベースのロギング、スタックロギング、コンテキスト統合などの高度な機能を備えた構造化ロギングの統一インターフェースを提供します。

## 特徴

- **高性能**：Logrus、Zap、標準 slog などの人気のロギングライブラリを上回る性能
- **構造化ロギング**：Go の標準`slog`パッケージをベースとした一貫した構造化出力
- **複数のロガータイプ**：
  - **標準ロガー**：設定可能なレベルを持つ高速で効率的なロギング
  - **タイマーロガー**：指定された間隔後にのみ出力するレート制限ロギング
  - **スタックロガー**：フィールド値をスタックに蓄積してトレース機能を提供
- **コンテキスト統合**：Go コンテキストへの第一級サポート
- **フィールドチェーン**：構造化フィールドを追加するためのフルエント API
- **設定可能な出力**：カスタム出力ライターのサポート
- **レベルベースロギング**：デバッグ、情報、警告、エラー、致命的エラーレベル
- **スレッドセーフ**：並行使用に安全

## インストール

```bash
go get github.com/yanun0323/logs
```

## クイックスタート

### 基本的な使用法

```go
package main

import (
    "github.com/yanun0323/logs"
)

func main() {
    // デフォルトロガーを使用
    logs.Info("アプリケーションが開始されました")
    logs.Error("何かが間違いました")

    // カスタムロガーを作成
    logger := logs.New(logs.LevelDebug)
    logger.WithField("user_id", "12345").Info("ユーザーがログインしました")

    // 複数のフィールドをチェーン
    logger.WithFields(map[string]interface{}{
        "method": "POST",
        "path":   "/api/users",
        "status": 201,
    }).Info("リクエストが完了しました")
}
```

### コンテキスト統合

```go
func handleRequest(ctx context.Context) {
    // ロガーをコンテキストにアタッチ
    logger := logs.New(logs.LevelInfo).WithField("request_id", "abc123")
    ctx = logger.Attach(ctx)

    // 後続の呼び出しチェーンで
    processData(ctx)
}

func processData(ctx context.Context) {
    // コンテキストからロガーを取得
    logger := logs.Get(ctx)
    logger.Info("データを処理中")
}
```

### タイマーロガー

タイマーロガーは、指定された間隔後にのみメッセージを出力することで、ログの氾濫を防ぐのに役立ちます：

```go
// 5秒に一度だけログ記録
timerLogger := logs.NewTimerLogger(5*time.Second, logs.LevelInfo)

for i := 0; i < 1000; i++ {
    // 1000回呼び出されても、数回しか出力されません
    timerLogger.Info("アイテムを処理中", i)
    time.Sleep(10 * time.Millisecond)
}
```

### スタックロガー

スタックロガーは指定されたフィールドキーの値を蓄積し、トレースのような出力を作成します：

```go
logger := logs.New(logs.LevelInfo)
stackLogger := logs.NewStackLogger(logger, "trace", "operation")

// スタックを構築
stackLogger = stackLogger.WithField("trace", "start")
stackLogger = stackLogger.WithField("operation", "validate")
stackLogger = stackLogger.WithField("trace", "middle")
stackLogger = stackLogger.WithField("operation", "process")

// 出力は表示されます：trace="start -> middle", operation="validate -> process"
stackLogger.Info("操作が完了しました")
```

## ロガータイプ

### 標準ロガー

- 高速で効率的なロギング
- Go の`slog`パッケージをベース
- 設定可能な出力先
- すべての標準ログレベルをサポート

### タイマーロガー

- レート制限によりログの氾濫を防止
- 設定可能な時間間隔
- スレッドセーフな実装
- 高頻度操作に有用

### スタックロガー

- フィールド値をスタックに蓄積
- 実行パスのトレースに最適
- 設定可能なスタックフィールドキー
- 呼び出し階層を維持

## パフォーマンスベンチマーク

他の人気の Go ロギングライブラリとのパフォーマンス比較：

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

**主要パフォーマンスハイライト：**

- **最高速**：Logrus より約 1.85 倍高速
- **メモリ効率**：Zap や Logrus と比較してメモリ割り当てが少ない
- **競争力**：追加機能を提供しながら標準 slog と同様のパフォーマンス

## API リファレンス

### ロガーインターフェース

```go
type Logger interface {
    // 基本的なロギングメソッド
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

    // 汎用ロギング
    Log(level Level, args ...interface{})
    Logf(level Level, format string, args ...interface{})

    // フィールド管理
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    WithContext(ctx context.Context) Logger

    // ユーティリティメソッド
    Copy() Logger
    Attach(ctx context.Context) context.Context
}
```

### ログレベル

```go
const (
    LevelDebug Level = -4
    LevelInfo  Level = 0
    LevelWarn  Level = 4
    LevelError Level = 8
    LevelFatal Level = 12
)
```

### 設定

```go
// デフォルトロガーを設定
logger := logs.New(logs.LevelInfo)
logs.SetDefault(logger)

// カスタム時間フォーマットを設定
logs.SetDefaultTimeFormat("2006-01-02 15:04:05")
```

## 依存関係

このライブラリの依存関係は最小限で、Go の標準ライブラリのみを使用します：

- **Go 1.21 以上** - 必要な Go バージョン
- **標準ライブラリのみ** - コア機能に外部依存関係は不要
- **オプションのテスト依存関係**：
  - `github.com/sirupsen/logrus` - パフォーマンスベンチマーク用
  - `go.uber.org/zap` - パフォーマンスベンチマーク用

このライブラリは軽量で依存関係のない設計でありながら、強力なロギング機能を提供します。

## コントリビューション

コントリビューションを歓迎します！お気軽に Pull Request を提出してください。

## ライセンス

このプロジェクトは MIT ライセンスの下でライセンスされています - 詳細については LICENSE ファイルをご覧ください。
