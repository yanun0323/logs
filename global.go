package logs

import (
	"context"

	"github.com/yanun0323/logs/internal"
)

var (
	initLogger = setupLogger()
	// defaultLogger is the default logger.
	defaultLogger = internal.NewValue(defaultLoggerWrapper{initLogger})
)

func setupLogger() Logger {
	return New(LevelInfo)
}

type defaultLoggerWrapper struct {
	logger Logger
}

// logKey is the key for the logger in the context.
type logKey struct{}

var logAttachKey = logKey{}

// Get gets the logger from context. if there's no logger in context, it will create a new logger with 'info' level.
func Get(ctx context.Context) Logger {
	val := ctx.Value(logAttachKey)
	if logger, ok := val.(Logger); ok {
		return logger
	}

	return Default()
}

// Default returns the default logger.
func Default() Logger {
	l, ok := defaultLogger.Load().(defaultLoggerWrapper)
	if !ok {
		l = defaultLoggerWrapper{logger: initLogger}
		defaultLogger.Store(l)
	}

	return l.logger
}

// SetDefault sets the default logger.
func SetDefault(logger Logger) {
	if logger != nil {
		defaultLogger.Store(defaultLoggerWrapper{logger: logger})
	}
}

// SetDefaultTimeFormat sets the default time format.
func SetDefaultTimeFormat(format string) {
	if len(format) != 0 {
		internal.SetDefaultTimeFormat(format)
	}
}

// Debug uses the default logger to log a message at the debug level.
func Debug(args ...any) {
	Default().Debug(args...)
}

// Debugf uses the default logger to log a message at the debug level.
func Debugf(format string, args ...any) {
	Default().Debugf(format, args...)
}

// Error uses the default logger to log a message at the error level.
func Error(args ...any) {
	Default().Error(args...)
}

// Errorf uses the default logger to log a message at the error level.
func Errorf(format string, args ...any) {
	Default().Errorf(format, args...)
}

// Fatal uses the default logger to log a message at the fatal level.
func Fatal(args ...any) {
	Default().Fatal(args...)
}

// Fatalf uses the default logger to log a message at the fatal level.
func Fatalf(format string, args ...any) {
	Default().Fatalf(format, args...)
}

// Info uses the default logger to log a message at the info level.
func Info(args ...any) {
	Default().Info(args...)
}

// Infof uses the default logger to log a message at the info level.
func Infof(format string, args ...any) {
	Default().Infof(format, args...)
}

// Warn uses the default logger to log a message at the warn level.
func Warn(args ...any) {
	Default().Warn(args...)
}

// Warnf uses the default logger to log a message at the warn level.
func Warnf(format string, args ...any) {
	Default().Warnf(format, args...)
}
