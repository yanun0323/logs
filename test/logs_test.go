package test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/yanun0323/logs"
)

func TestSetDefault(t *testing.T) {
	log := logs.NewTraceLogger(logs.LevelInfo, logs.KeyFunc)
	logs.SetDefault(log)
	log.Info("Test")

	log2 := logs.New(logs.LevelInfo)
	logs.SetDefault(log2)
	log2.Info("Test")
}

func TestGet(t *testing.T) {
	log := logs.Get(context.Background())
	log.Info("Test")
}

func TestLogOutput(t *testing.T) {
	w2 := logs.FileOutput(".", "dir_dot")
	w3 := logs.FileOutput("", "dir_empty")
	w4 := logs.FileOutput("hello", "dir_wrong")

	log1 := logs.New(logs.LevelInfo)
	log2 := logs.New(logs.LevelInfo, &logs.Option{Output: w2})
	log3 := logs.New(logs.LevelInfo, &logs.Option{Output: w3})
	log4 := logs.New(logs.LevelInfo, &logs.Option{Output: w4})

	t.Logf("log1 = %p, log2 = %p, log3 = %p, log4 = %p", log1, log2, log3, log4)
	log1.Info("info")
	log2.Info("info")
	log3.Info("info")
	log4.Info("info")

	if err := w2.Remove(); err != nil {
		t.Errorf("remove w2 failed: %v", err)
	}

	if err := w3.Remove(); err != nil {
		t.Errorf("remove w3 failed: %v", err)
	}

	if err := w4.Remove(); err != nil {
		t.Errorf("remove w4 failed: %v", err)
	}
}

func TestLogs(t *testing.T) {
	log1 := logs.New(logs.LevelInfo)
	log2 := logs.New(logs.LevelInfo)

	t.Logf("log1 = %p, log2 = %p", log1, log2)
	log1.Info("info")
	log2.Info("info")
}

func TestMap(t *testing.T) {
	log1 := logs.New(logs.LevelInfo)
	log1.With("test", map[string]any{"test": true}).Info("access")
}

func TestWith(t *testing.T) {
	log := logs.New(logs.LevelInfo).With("hello", "foo....")

	log.Info("hello field info...")
	log.With("user_id", "i'm user").With("info_id", "i'm order").Info("with user id")
}

func TestWiths(t *testing.T) {
	log := logs.New(logs.LevelInfo).With(map[string]any{
		"foo": 123,
		"bar": "456",
	})

	log.Info("info...")
	log.With("user_id", "i'm user").Info("with user id")
}

func TestFatal(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		logs.New(logs.LevelInfo).Fatal("fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogs_Fatal")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestLevel(t *testing.T) {
	writer := &bytes.Buffer{}
	log := logs.New(logs.LevelWarn, &logs.Option{Output: writer})
	log.Info("LEVEL info")
	if writer.Len() != 0 {
		t.Errorf("writer len is not 0")
	}

	writer.Reset()

	log.With("user_id", "i'm user").
		Warn("LEVEL with user id")

	if writer.Len() == 0 {
		t.Errorf("writer len is 0")
	}
}

func TestExample(t *testing.T) {
	l := logs.NewTraceLogger(logs.LevelDebug, "func", &logs.Option{
		// Output: logs.FileOutput(".", "test_example"),
		// Format: logs.FormatText,
	})
	ctx := context.TODO()

	println()
	l.Debug("debug message")
	l.With("fields", "val").
		With(logs.KeyErr, nil).
		With(logs.KeyCtx, ctx).
		With(logs.KeyFunc, "testFunc").
		Info("info message with fields")

	l.With("func", "F0").
		With("func", "F1", "func", "F2").
		With("func", "F3").
		Warn("warn message with func trace")

	l.Error("error message")
	l.Fatal("fatal message")
	println()
}
