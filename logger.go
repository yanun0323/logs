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

type logger slog.Logger

// New creates a new basic logger with the given level and outputs.
//
// If option is not provided, the logger will write to the os.Stdout with console format.
func New(level Level, option ...*Option) Logger {
	if len(option) != 0 {
		return (*logger)(slog.New(option[0].createLoggerHandler(level)))
	}

	return (*logger)(slog.New(defaultOption.createLoggerHandler(level)))
}

func (l *logger) Copy() Logger {
	return (*logger)((*slog.Logger)(l))
}

func (l *logger) With(args ...any) Logger {
	if len(args) == 0 {
		return l
	}

	return (*logger)((*slog.Logger)(l).With(args...))
}

func (l *logger) WithError(err error) Logger {
	return l.With(KeyErr, err)
}

func (l *logger) WithCtx(ctx context.Context) Logger {
	return l.With(KeyCtx, ctx)
}

func (l *logger) WithFunc(function string) Logger {
	return l.With(KeyFunc, function)
}

func (l *logger) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, logAttachKey, l)
}

func (l *logger) Log(level Level, args ...any) {
	slogLevel := slog.Level(level)

	switch len(args) {
	case 0:
		(*slog.Logger)(l).Log(bgCtx, slogLevel, "")
	case 1:
		if str, ok := args[0].(string); ok {
			(*slog.Logger)(l).Log(bgCtx, slogLevel, str)
		} else {
			(*slog.Logger)(l).Log(bgCtx, slogLevel, internal.ValueToString(args[0]))
		}
	case 2:
		(*slog.Logger)(l).Log(bgCtx, slogLevel, fmt.Sprint(args[0], " ", args[1]))
	default:
		(*slog.Logger)(l).Log(bgCtx, slogLevel, fmt.Sprint(args...))
	}
}

func (l *logger) Logf(level Level, format string, args ...any) {
	slogLevel := slog.Level(level)

	if len(args) == 0 {
		(*slog.Logger)(l).Log(bgCtx, slogLevel, format)
		return
	}

	(*slog.Logger)(l).Log(bgCtx, slogLevel, fmt.Sprintf(format, args...))
}

func (l *logger) Debug(args ...any) {
	l.Log(LevelDebug, args...)
}

func (l *logger) Debugf(format string, args ...any) {
	l.Logf(LevelDebug, format, args...)
}

func (l *logger) Info(args ...any) {
	switch len(args) {
	case 0:
		(*slog.Logger)(l).Log(bgCtx, slog.LevelInfo, "")
	case 1:
		if str, ok := args[0].(string); ok {
			(*slog.Logger)(l).Log(bgCtx, slog.LevelInfo, str)
		} else {
			(*slog.Logger)(l).Log(bgCtx, slog.LevelInfo, fmt.Sprint(args[0]))
		}
	case 2:
		(*slog.Logger)(l).Log(bgCtx, slog.LevelInfo, fmt.Sprint(args[0], " ", args[1]))
	default:
		(*slog.Logger)(l).Log(bgCtx, slog.LevelInfo, fmt.Sprint(args...))
	}
}

func (l *logger) Infof(format string, args ...any) {
	if len(args) == 0 {
		(*slog.Logger)(l).Log(bgCtx, slog.LevelInfo, format)
		return
	}
	(*slog.Logger)(l).Log(bgCtx, slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *logger) Warn(args ...any) {
	l.Log(LevelWarn, args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.Logf(LevelWarn, format, args...)
}

func (l *logger) Error(args ...any) {
	switch len(args) {
	case 0:
		(*slog.Logger)(l).Log(bgCtx, slog.LevelError, "")
	case 1:
		if str, ok := args[0].(string); ok {
			(*slog.Logger)(l).Log(bgCtx, slog.LevelError, str)
		} else {
			(*slog.Logger)(l).Log(bgCtx, slog.LevelError, fmt.Sprint(args[0]))
		}
	case 2:
		(*slog.Logger)(l).Log(bgCtx, slog.LevelError, fmt.Sprint(args[0], " ", args[1]))
	default:
		(*slog.Logger)(l).Log(bgCtx, slog.LevelError, fmt.Sprint(args...))
	}
}

func (l *logger) Errorf(format string, args ...any) {
	if len(args) == 0 {
		(*slog.Logger)(l).Log(bgCtx, slog.LevelError, format)
		return
	}
	(*slog.Logger)(l).Log(bgCtx, slog.LevelError, fmt.Sprintf(format, args...))
}

func (l *logger) Fatal(args ...any) {
	l.Log(LevelFatal, args...)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.Logf(LevelFatal, format, args...)
	os.Exit(1)
}
