package logs

import (
	"bytes"
	"context"
	"fmt"

	"github.com/yanun0323/logs/internal/buffer"
)

type traceLogger struct {
	Logger

	keyword string
	stack   *bytes.Buffer
}

// NewTraceLogger creates a new trace logger with the given level and trace field key.
//
// It will accumulate the values of the specified field keys into a stack,
// and the stack will be outputted when the logger is called.
//
// If option is not provided, the logger will write to the os.Stdout with console format.
func NewTraceLogger(level Level, traceFieldKeyword string, option ...*Option) Logger {
	if traceFieldKeyword == "" {
		return New(level, option...)
	}

	return &traceLogger{
		keyword: traceFieldKeyword,
		Logger:  New(level, option...),
		stack:   &bytes.Buffer{},
	}
}

func (l *traceLogger) clone() *traceLogger {
	buf := l.stack.Bytes()
	stack := buffer.Pool.Get().(*bytes.Buffer)
	stack.Reset()
	stack.Grow(len(buf) + 256)
	stack.Write(buf)

	return &traceLogger{
		keyword: l.keyword,
		Logger:  l.Logger.Copy(),
		stack:   stack,
	}
}

func (l *traceLogger) Copy() Logger {
	return l.clone()
}

func (l *traceLogger) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey{}, l)
}

const (
	_stackSeparator = " -> "
)

func (l *traceLogger) WithField(key string, value any) Logger {
	if key == l.keyword {
		buf := l.stack.Bytes()
		str := fmt.Sprintf("%v", value)
		stack := buffer.Pool.Get().(*bytes.Buffer)
		stack.Reset()
		stack.Grow(len(buf) + len(str) + len(_stackSeparator))
		stack.Write(buf)
		if stack.Len() != 0 {
			stack.WriteString(_stackSeparator)
		}
		stack.WriteString(str)

		return &traceLogger{
			keyword: l.keyword,
			Logger:  l.Logger,
			stack:   stack,
		}
	} else {
		return &traceLogger{
			keyword: l.keyword,
			Logger:  l.Logger.WithField(key, value),
			stack:   l.stack,
		}
	}
}

func (l *traceLogger) WithFields(fields map[string]any) Logger {
	if len(fields) == 0 {
		return l
	}

	var (
		hasStackFields bool
		stackFields    = make(map[string][]any)
		normalFields   = make(map[string]any, len(fields))
	)

	for k, v := range fields {
		if k == l.keyword {
			stackFields[k] = append(stackFields[k], v)
			hasStackFields = true
		} else {
			normalFields[k] = v
		}
	}

	logger := l.Logger
	stack := l.stack

	if len(normalFields) != 0 {
		logger = logger.WithFields(normalFields)
	}

	if hasStackFields {
		buf := l.stack.Bytes()
		for _, val := range stackFields {
			stack = buffer.Pool.Get().(*bytes.Buffer)
			stack.Reset()
			stack.Grow(len(buf) + (len(_stackSeparator)+256)*len(val))
			stack.Write(buf)
			for _, v := range val {
				if stack.Len() != 0 {
					stack.WriteString(_stackSeparator)
				}
				stack.WriteString(fmt.Sprintf("%v", v))
			}
		}
	}

	return &traceLogger{
		keyword: l.keyword,
		Logger:  logger,
		stack:   stack,
	}
}

func (l *traceLogger) fieldsToAttach() map[string]any {
	if l.stack.Len() == 0 {
		return nil
	}

	return map[string]any{
		l.keyword: l.stack.String(),
	}
}

func (l *traceLogger) withFieldsIfNeeded() Logger {
	fields := l.fieldsToAttach()
	if len(fields) == 0 {
		return l.Logger
	}
	return l.Logger.WithFields(fields)
}

func (l *traceLogger) Log(level Level, args ...any) {
	l.withFieldsIfNeeded().Log(level, args...)
}

func (l *traceLogger) Logf(level Level, format string, args ...any) {
	l.withFieldsIfNeeded().Logf(level, format, args...)
}

func (l *traceLogger) Debug(args ...any) {
	l.withFieldsIfNeeded().Debug(args...)
}

func (l *traceLogger) Debugf(format string, args ...any) {
	l.withFieldsIfNeeded().Debugf(format, args...)
}

func (l *traceLogger) Info(args ...any) {
	l.withFieldsIfNeeded().Info(args...)
}

func (l *traceLogger) Infof(format string, args ...any) {
	l.withFieldsIfNeeded().Infof(format, args...)
}

func (l *traceLogger) Warn(args ...any) {
	l.withFieldsIfNeeded().Warn(args...)
}

func (l *traceLogger) Warnf(format string, args ...any) {
	l.withFieldsIfNeeded().Warnf(format, args...)
}

func (l *traceLogger) Error(args ...any) {
	l.withFieldsIfNeeded().Error(args...)
}

func (l *traceLogger) Errorf(format string, args ...any) {
	l.withFieldsIfNeeded().Errorf(format, args...)
}

func (l *traceLogger) Fatal(args ...any) {
	l.withFieldsIfNeeded().Fatal(args...)
}

func (l *traceLogger) Fatalf(format string, args ...any) {
	l.withFieldsIfNeeded().Fatalf(format, args...)
}

func (l *traceLogger) WithError(err error) Logger {
	return &traceLogger{
		keyword: l.keyword,
		Logger:  l.Logger.WithError(err),
		stack:   l.stack,
	}
}

func (l *traceLogger) WithContext(ctx context.Context) Logger {
	return &traceLogger{
		keyword: l.keyword,
		Logger:  l.Logger.WithContext(ctx),
		stack:   l.stack,
	}
}

func (l *traceLogger) WithFunc(function string) Logger {
	return &traceLogger{
		keyword: l.keyword,
		Logger:  l.Logger.WithFunc(function),
		stack:   l.stack,
	}
}
