package logs

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestTickerLogger(t *testing.T) {
	writer := &bytes.Buffer{}
	timer := NewTickerLogger(LevelDebug, time.Hour, &Option{Output: writer})
	timer.Debug("debug")
	timer.Info("info")
	timer.Warn("warn")
	timer.Error("error")

	if line := len(strings.Split(strings.TrimSpace(writer.String()), "\n")); line != 1 {
		t.Errorf("Expected one line, but got %d lines: %s", line, writer.String())
	}

	writer = &bytes.Buffer{}
	timer = NewTickerLogger(LevelDebug, time.Microsecond, &Option{Output: writer})
	time.Sleep(time.Microsecond)
	timer.Debug("debug")
	timer.Info("info")
	timer.Warn("warn")
	timer.Error("error")

	if writer.Len() == 0 {
		t.Error("Expected output, but got none")
	}
}

func TestTickerLoggerInterval(t *testing.T) {
	writer := &bytes.Buffer{}
	timer := NewTickerLogger(LevelDebug, time.Second, &Option{Output: writer})

	for i := 0; i < 10; i++ {
		timer.Debug("debug")
		time.Sleep(100 * time.Millisecond)
	}

	result := writer.String()
	if len(result) == 0 {
		t.Error("Expected at least one log output")
	}

	// Should have maximum 2 log entries in 1 second with 100ms intervals
	// (one immediate, one after 1 second)
	lines := bytes.Count(writer.Bytes(), []byte("\n"))
	if lines > 3 { // Allow some tolerance
		t.Errorf("Expected at most 3 log lines, got %d", lines)
	}
}
