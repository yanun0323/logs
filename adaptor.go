package logs

import (
	"context"
)

// Logger is the interface that wraps the basic methods of a logger.
type Logger interface {
	// Copy duplicates the logger.
	Copy() Logger

	// Attach attaches the logger into the context.
	Attach(ctx context.Context) context.Context

	// WithField copies the logger and adds a single field to the Logger.
	WithField(key string, value any) Logger
	// WithFields copies the logger and adds a map of fields to the Logger.
	WithFields(fields map[string]any) Logger
	// WithError copies the logger and adds an error as single field (using the key defined in ErrorKey) to the Logger.
	WithError(err error) Logger
	// WithFunc copies the logger and adds a function to the Logger.
	WithFunc(function string) Logger

	// WithContext copies the logger and adds a context to the Logger.
	WithContext(ctx context.Context) Logger

	// Log will log a message at the level given as parameter.
	//
	// Warning: using Log at Fatal level will not respectively Exit.
	// For this behavior Entry.Fatal should be used instead.
	Log(level Level, args ...any)

	// Logf will log a message at the level given as parameter.
	//
	// Warning: using Log at Fatal level will not respectively Exit.
	// For this behavior Entry.Fatal should be used instead.
	Logf(level Level, format string, args ...any)

	// Debug will log a message at the debug level.
	Debug(args ...any)

	// Debugf will log a message at the debug level.
	Debugf(format string, args ...any)

	// Info will log a message at the info level.
	Info(args ...any)

	// Infof will log a message at the info level.
	Infof(format string, args ...any)

	// Warn will log a message at the warn level.
	Warn(args ...any)

	// Warnf will log a message at the warn level.
	Warnf(format string, args ...any)

	// Error will log a message at the error level.
	Error(args ...any)

	// Errorf will log a message at the error level.
	Errorf(format string, args ...any)

	// Fatal will log a message at the fatal level.
	Fatal(args ...any)

	// Fatalf will log a message at the fatal level.
	Fatalf(format string, args ...any)
}
