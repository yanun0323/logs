package logs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type loggerNew slog.Logger

// New creates a new basic logger with the given level and outputs.
//
// If option is not provided, the logger will write to the os.Stdout with console format.
func New(level Level, option ...*Option) Logger {
	if len(option) != 0 {
		return (*loggerNew)(slog.New(option[0].createLoggerHandler(level)))
	}

	return (*loggerNew)(slog.New(defaultOption.createLoggerHandler(level)))
}

func (l loggerNew) clone() *loggerNew {
	return (*loggerNew)((*slog.Logger)(&l))
}

func (l *loggerNew) Copy() Logger {
	return l.clone()
}

func (l *loggerNew) withField(key string, value any) *loggerNew {
	return (*loggerNew)((*slog.Logger)(l).With(key, value))
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

	return (*loggerNew)((*slog.Logger)(l).With(attrs...))
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
	if len(args) == 0 {
		(*slog.Logger)(l).Log(context.Background(), slog.Level(level), "")
		return
	}
	if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			(*slog.Logger)(l).Log(context.Background(), slog.Level(level), str)
			return
		}
	}
	(*slog.Logger)(l).Log(context.Background(), slog.Level(level), fmt.Sprint(args...))
}

func (l *loggerNew) Logf(level Level, format string, args ...any) {
	if len(args) == 0 {
		(*slog.Logger)(l).Log(context.Background(), slog.Level(level), format)
	} else {
		(*slog.Logger)(l).Log(context.Background(), slog.Level(level), fmt.Sprintf(format, args...))
	}
}

func (l *loggerNew) Debug(args ...any) {
	l.Log(LevelDebug, args...)
}

func (l *loggerNew) Debugf(format string, args ...any) {
	l.Logf(LevelDebug, format, args...)
}

func (l *loggerNew) Info(args ...any) {
	l.Log(LevelInfo, args...)
}

func (l *loggerNew) Infof(format string, args ...any) {
	l.Logf(LevelInfo, format, args...)
}

func (l *loggerNew) Warn(args ...any) {
	l.Log(LevelWarn, args...)
}

func (l *loggerNew) Warnf(format string, args ...any) {
	l.Logf(LevelWarn, format, args...)
}

func (l *loggerNew) Error(args ...any) {
	l.Log(LevelError, args...)
}

func (l *loggerNew) Errorf(format string, args ...any) {
	l.Logf(LevelError, format, args...)
}

func (l *loggerNew) Fatal(args ...any) {
	l.Log(LevelFatal, args...)
	os.Exit(1)
}

func (l *loggerNew) Fatalf(format string, args ...any) {
	l.Logf(LevelFatal, format, args...)
	os.Exit(1)
}
