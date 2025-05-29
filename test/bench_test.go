package test

import (
	"log/slog"
	"testing"
	"time"

	zerolog "github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/yanun0323/logs"
	"github.com/yanun0323/logs/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
			l.WithFields("key", "value", "key2", 123.456).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.WithFields("key", "value", "key2", 123.456).Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}
func BenchmarkLogsTicker(b *testing.B) {
	writer := switchableWriter(".", "logger.ticker.log")

	l := logs.NewTickerLogger(logs.LevelInfo, time.Second, &logs.Option{Output: writer})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.WithFields("key", "value", "key2", 123.456).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.WithFields("key", "value", "key2", 123.456).Info("test")
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
			l.WithFields("key", "value", "key2", 123.456).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.WithFields("key", "value", "key2", 123.456).Info("test")
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
			l.With("key", "value", "key2", 123.456).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.With("key", "value", "key2", 123.456).Info("test")
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
			l.With("key", "value", "key2", 123.456).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.With("key", "value", "key2", 123.456).Info("test")
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
			l.With("key", "value", "key2", 123.456).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.With("key", "value", "key2", 123.456).Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerS failed: %v", err)
		}
	})
}

func BenchmarkZap(b *testing.B) {
	writer := switchableWriter(".", "zap.log")
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeLevel = zapcore.LowercaseColorLevelEncoder

	l := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(conf),
		writer,
		zap.NewAtomicLevelAt(zap.InfoLevel),
	))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.With(zap.Any("key", "value"), zap.Any("key2", 123.456)).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.With(zap.Any("key", "value"), zap.Any("key2", 123.456)).Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerZ failed: %v", err)
		}
	})
}

func BenchmarkLogrus(b *testing.B) {
	writer := switchableWriter(".", "logrus.log")

	l := logrus.New()
	l.Out = writer
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info("test")
			l.WithFields(map[string]any{
				"key":  "value",
				"key2": 123.456,
			}).Info("test")
		}
	})

	for b.Loop() {
		l.Info("test")
		l.WithFields(map[string]any{
			"key":  "value",
			"key2": 123.456,
		}).Info("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}

func BenchmarkZeroLog(b *testing.B) {
	writer := switchableWriter(".", "zerolog.log")

	l := zerolog.Logger.Output(writer)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.Info().Msg("test")
			ll := l.With().Any("key", "value").Any("key2", 123.456).Logger()
			ll.Info().Msg("test")
		}
	})

	for b.Loop() {
		l.Info().Msg("test")
		ll := l.With().Any("key", "value").Any("key2", 123.456).Logger()
		ll.Info().Msg("test")
	}

	b.Cleanup(func() {
		if err := writer.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}
