package logs

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/yanun0323/logs/internal"
)

var (
	bgCtx = context.Background()
)

type loggerNew struct {
	*slog.Logger
}

// New creates a new basic logger with the given level and outputs.
//
// If option is not provided, the logger will write to the os.Stdout with console format.
func New(level Level, option ...*Option) Logger {
	if len(option) != 0 {
		return &loggerNew{Logger: slog.New(option[0].createLoggerHandler(level))}
	}

	return &loggerNew{Logger: slog.New(defaultOption.createLoggerHandler(level))}
}

func (l *loggerNew) clone() *loggerNew {
	return &loggerNew{Logger: l.Logger}
}

func (l *loggerNew) Copy() Logger {
	return l.clone()
}

func (l *loggerNew) withField(key string, value any) *loggerNew {
	return &loggerNew{Logger: l.Logger.With(key, value)}
}

func (l *loggerNew) WithField(key string, value any) Logger {
	return l.withField(key, value)
}

func (l *loggerNew) WithFields(fields map[string]any) Logger {
	if len(fields) == 0 {
		return l
	}

	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}

	return &loggerNew{Logger: l.Logger.With(attrs...)}
}

func (l *loggerNew) WithError(err error) Logger {
	return l.WithField(FieldKeyError, err)
}

func (l *loggerNew) WithContext(ctx context.Context) Logger {
	return l.WithField(FieldKeyContext, ctx)
}

func (l *loggerNew) WithFunc(function string) Logger {
	return l.WithField(FieldKeyFunc, function)
}

func (l *loggerNew) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey{}, l)
}

func (l *loggerNew) Log(level Level, args ...any) {
	slogLevel := slog.Level(level)

	switch len(args) {
	case 0:
		l.Logger.Log(bgCtx, slogLevel, "")
	case 1:
		if str, ok := args[0].(string); ok {
			l.Logger.Log(bgCtx, slogLevel, str)
		} else {
			// 直接使用 fmt.Sprint 但只針對單個值
			l.Logger.Log(bgCtx, slogLevel, internal.ValueToString(args[0]))
		}
	case 2:
		l.Logger.Log(bgCtx, slogLevel, fmt.Sprint(args[0], " ", args[1]))
	default:
		l.Logger.Log(bgCtx, slogLevel, fmt.Sprint(args...))
	}
}

func (l *loggerNew) Logf(level Level, format string, args ...any) {
	slogLevel := slog.Level(level)

	if len(args) == 0 {
		l.Logger.Log(bgCtx, slogLevel, format)
		return
	}

	l.Logger.Log(bgCtx, slogLevel, fmt.Sprintf(format, args...))
}

func (l *loggerNew) Debug(args ...any) {
	l.Log(LevelDebug, args...)
}

func (l *loggerNew) Debugf(format string, args ...any) {
	l.Logf(LevelDebug, format, args...)
}

func (l *loggerNew) Info(args ...any) {
	switch len(args) {
	case 0:
		l.Logger.Log(bgCtx, slog.LevelInfo, "")
	case 1:
		if str, ok := args[0].(string); ok {
			l.Logger.Log(bgCtx, slog.LevelInfo, str)
		} else {
			l.Logger.Log(bgCtx, slog.LevelInfo, fmt.Sprint(args[0]))
		}
	case 2:
		l.Logger.Log(bgCtx, slog.LevelInfo, fmt.Sprint(args[0], " ", args[1]))
	default:
		l.Logger.Log(bgCtx, slog.LevelInfo, fmt.Sprint(args...))
	}
}

func (l *loggerNew) Infof(format string, args ...any) {
	if len(args) == 0 {
		l.Logger.Log(bgCtx, slog.LevelInfo, format)
		return
	}
	l.Logger.Log(bgCtx, slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *loggerNew) Warn(args ...any) {
	l.Log(LevelWarn, args...)
}

func (l *loggerNew) Warnf(format string, args ...any) {
	l.Logf(LevelWarn, format, args...)
}

func (l *loggerNew) Error(args ...any) {
	switch len(args) {
	case 0:
		l.Logger.Log(bgCtx, slog.LevelError, "")
	case 1:
		if str, ok := args[0].(string); ok {
			l.Logger.Log(bgCtx, slog.LevelError, str)
		} else {
			l.Logger.Log(bgCtx, slog.LevelError, fmt.Sprint(args[0]))
		}
	case 2:
		l.Logger.Log(bgCtx, slog.LevelError, fmt.Sprint(args[0], " ", args[1]))
	default:
		l.Logger.Log(bgCtx, slog.LevelError, fmt.Sprint(args...))
	}
}

func (l *loggerNew) Errorf(format string, args ...any) {
	if len(args) == 0 {
		l.Logger.Log(bgCtx, slog.LevelError, format)
		return
	}
	l.Logger.Log(bgCtx, slog.LevelError, fmt.Sprintf(format, args...))
}

func (l *loggerNew) Fatal(args ...any) {
	l.Log(LevelFatal, args...)
	os.Exit(1)
}

func (l *loggerNew) Fatalf(format string, args ...any) {
	l.Logf(LevelFatal, format, args...)
	os.Exit(1)
}
