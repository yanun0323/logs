package logs

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func TestTraceLogger(t *testing.T) {
	writer := &bytes.Buffer{}
	trace := NewTraceLogger(LevelDebug, "func", writer)

	trace = trace.WithField("func", "func_1")
	trace = trace.WithField("func", "func_2")
	trace = trace.WithField("func", "func_3")

	trace.Debug("debug")

	result := writer.String()
	if !strings.Contains(result, "func_1 -> func_2 -> func_3") {
		t.Errorf("Expected trace to contain function chain, got: %s", result)
	}
}

func TestTraceLoggerContext(t *testing.T) {
	writer := &bytes.Buffer{}
	trace := NewTraceLogger(LevelDebug, "func", writer)

	trace = trace.
		WithField("func", "main").
		WithField("keyword", "A").
		WithField("single", "A")

	ctx := trace.Attach(context.Background())

	trace = Get(ctx)
	trace = trace.
		WithField("func", "sub").
		WithField("keyword", "B").
		WithField("single", "B")

	trace.Info("info")

	all, err := io.ReadAll(writer)
	if err != nil {
		t.Fatalf("read all failed: %v", err)
	}

	t.Log(string(all))
}
