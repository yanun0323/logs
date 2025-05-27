package test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/yanun0323/logs"
	"github.com/yanun0323/logs/internal"
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
			l.WithField("key", "value").Info("test")
		}
	})

	for i := 0; i < b.N; i++ {
		l.Info("test")
		l.WithField("key", "value").Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}
func BenchmarkLogsTicker(b *testing.B) {
	writer := switchableWriter(".", "logger.ticker.log")

	l := logs.NewTickerLogger(time.Second, logs.LevelInfo, &logs.Option{Output: writer})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.WithField("key", "value").Info("test")
		}
	})

	for i := 0; i < b.N; i++ {
		l.Info("test")
		l.WithField("key", "value").Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}

func BenchmarkLogsTrace(b *testing.B) {
	writer := switchableWriter(".", "logger.trace.log")

	l := logs.NewTraceLogger(logs.LevelInfo, "key", &logs.Option{Output: writer})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.WithField("key", "value").Info("test")
		}
	})

	for i := 0; i < b.N; i++ {
		l.Info("test")
		l.WithField("key", "value").Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}

func BenchmarkSlogWithTextHandler(b *testing.B) {
	writer := switchableWriter(".", "slog.log")

	l := slog.New(slog.NewTextHandler(writer, nil))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.With("key", "value").Info("test")
		}
	})

	for i := 0; i < b.N; i++ {
		l.Info("test")
		l.With("key", "value").Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerS failed: %v", err)
		}
	})
}

func BenchmarkSlogWithJSONHandler(b *testing.B) {
	writer := switchableWriter(".", "slog.json.log")

	l := slog.New(slog.NewJSONHandler(writer, nil))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.With("key", "value").Info("test")
		}
	})

	for i := 0; i < b.N; i++ {
		l.Info("test")
		l.With("key", "value").Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerS failed: %v", err)
		}
	})
}

func BenchmarkSlogLogsHandler(b *testing.B) {
	writer := switchableWriter(".", "slog.logs_handler.log")

	l := slog.New(internal.NewLoggerHandler(writer, int8(logs.LevelInfo)))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.With("key", "value").Info("test")
		}
	})

	for i := 0; i < b.N; i++ {
		l.Info("test")
		l.With("key", "value").Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerS failed: %v", err)
		}
	})
}

// func BenchmarkZap(b *testing.B) {
// 	writer := switchableWriter(".", "zap.log")
// 	conf := zap.NewProductionEncoderConfig()
// 	conf.EncodeLevel = zapcore.LowercaseColorLevelEncoder

// 	l := zap.New(zapcore.NewCore(
// 		zapcore.NewConsoleEncoder(conf),
// 		writer,
// 		zap.NewAtomicLevelAt(zap.InfoLevel),
// 	))
// 	b.RunParallel(func(p *testing.PB) {
// 		for p.Next() {
// 			l.Info("test")
// 			l.With(zap.Any("key", "value")).Info("test")
// 		}
// 	})

// 	for i := 0; i < b.N; i++ {
// 		l.Info("test")
// 		l.With(zap.Any("key", "value")).Info("test")
// 	}

// 	b.Cleanup(func() {
// 		if err := writer.Remove(); err != nil {
// 			b.Fatalf("remove writerZ failed: %v", err)
// 		}
// 	})
// }

// func BenchmarkLogrus(b *testing.B) {
// 	writer := switchableWriter(".", "logrus.log")

// 	l := logrus.New()
// 	l.Out = writer
// 	b.RunParallel(func(p *testing.PB) {
// 		for p.Next() {
// 			l.Info("test")
// 			l.WithField("key", "value").Info("test")
// 		}
// 	})

// 	for i := 0; i < b.N; i++ {
// 		l.Info("test")
// 		l.WithField("key", "value").Info("test")
// 	}

// 	b.Cleanup(func() {
// 		if err := writer.Remove(); err != nil {
// 			b.Fatalf("remove writerL failed: %v", err)
// 		}
// 	})
// }
