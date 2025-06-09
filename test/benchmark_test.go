package test

import (
	"testing"

	"github.com/yanun0323/logs"
)

func switchableWriter(relativeDir string, filename string) logs.Writer {
	return logs.FileOutput(relativeDir, filename)
}

func BenchmarkLogsBasic(b *testing.B) {
	writer := switchableWriter(".", "logger.basic.log")

	l := logs.New(logs.LevelInfo, &logs.Option{Output: writer})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.With("key", "value", "key2", 123.456).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.With("key", "value", "key2", 123.456).Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}
