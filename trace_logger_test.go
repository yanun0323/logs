package logs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestTraceLogger(t *testing.T) {
	writer := &bytes.Buffer{}
	trace := NewTraceLogger(LevelDebug, "func", &Option{Output: writer})

	trace = trace.With("func", "func_1").With("hello", "world")
	trace = trace.With("func", "func_2")
	trace = trace.With("func", "func_3")

	trace.Debug("debug")

	result := writer.String()
	t.Log(result)
	if !strings.Contains(result, "func_1 -> func_2 -> func_3") {
		t.Errorf("Expected trace to contain function chain, got: %s", result)
	}
}

func TestTraceLoggerContext(t *testing.T) {
	writer := &bytes.Buffer{}
	trace := NewTraceLogger(LevelDebug, "func", &Option{Output: writer})

	trace = trace.
		With("func", "main").
		With("keyword", "A").
		With("single", "A")

	ctx := trace.Attach(context.Background())

	trace = Get(ctx)
	trace = trace.
		With("func", "sub").
		With("keyword", "B").
		With("single", "B")

	trace.Info("info")

	all, err := io.ReadAll(writer)
	if err != nil {
		t.Fatalf("read all failed: %v", err)
	}

	t.Log(string(all))
}

func BenchmarkTraceLoggerWith(b *testing.B) {
	trace := NewTraceLogger(LevelInfo, "trace", &Option{Output: EmptyOutput})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			trace.With("trace", "function_name").Info("test message")
		}
	})
}

func BenchmarkTraceLoggerWithMultipleFields(b *testing.B) {
	trace := NewTraceLogger(LevelInfo, "trace", &Option{Output: EmptyOutput})

	fields := map[string]any{
		"trace":   "function_name",
		"user":    12345,
		"action":  "login",
		"success": true,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			trace.With(fields).Info("test message")
		}
	})
}

func BenchmarkTraceLoggerStackBuilding(b *testing.B) {
	trace := NewTraceLogger(LevelInfo, "func", &Option{Output: EmptyOutput})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := trace
		for j := 0; j < 10; j++ {
			logger = logger.With("func", fmt.Sprintf("func_%d", j))
		}
		logger.Info("final message")
	}
}

// 測試深度嵌套的效能
func BenchmarkTraceLoggerDeepNesting(b *testing.B) {
	trace := NewTraceLogger(LevelInfo, "trace", &Option{Output: EmptyOutput})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger := trace
			// 模擬深度呼叫堆疊
			for i := 0; i < 20; i++ {
				logger = logger.With("trace", fmt.Sprintf("depth_%d", i))
			}
			logger.Info("deep call")
		}
	})
}

// 測試混合字段的效能
func BenchmarkTraceLoggerMixedFields(b *testing.B) {
	trace := NewTraceLogger(LevelInfo, "trace", &Option{Output: EmptyOutput})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			trace.
				With("trace", "service_call").
				With("user_id", 12345).
				With("trace", "database_query").
				With("duration", 150.5).
				With("trace", "result_processing").
				Info("request completed")
		}
	})
}

// 測試記憶體分配情況
func BenchmarkTraceLoggerMemoryAllocation(b *testing.B) {
	trace := NewTraceLogger(LevelInfo, "func", &Option{Output: EmptyOutput})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger := trace.With("func", "test_function")
		logger.Info("allocation test")
	}
}
