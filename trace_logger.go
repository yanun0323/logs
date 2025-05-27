package logs

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
)

// stringBuilderPool 重用 strings.Builder，減少記憶體分配
var stringBuilderPool = sync.Pool{
	New: func() any {
		return &strings.Builder{}
	},
}

type traceLogger struct {
	Logger

	keyword string
	stackMu sync.RWMutex
	stack   map[string][]any
}

// NewTraceLogger creates a new trace logger with the given level and trace field key.
//
// It will accumulate the values of the specified field keys into a stack,
// and the stack will be outputted when the logger is called.
//
// If outputs is not provided, the logger will write to the os.Stdout.
func NewTraceLogger(level Level, traceFieldKeyword string, outputs ...io.Writer) Logger {
	if traceFieldKeyword == "" {
		return New(level, outputs...)
	}

	return &traceLogger{
		keyword: traceFieldKeyword,
		Logger:  New(level, outputs...),
		stack:   make(map[string][]any),
	}
}

// 優化：使用更高效的 stack 複製，避免不必要的分配
func (l *traceLogger) copyStack() map[string][]any {
	l.stackMu.RLock()
	defer l.stackMu.RUnlock()

	if len(l.stack) == 0 {
		return make(map[string][]any)
	}

	stackCopied := make(map[string][]any, len(l.stack))
	for k, v := range l.stack {
		if len(v) > 0 {
			stackCopied[k] = make([]any, len(v))
			copy(stackCopied[k], v)
		}
	}

	return stackCopied
}

func (l *traceLogger) clone() *traceLogger {
	return &traceLogger{
		keyword: l.keyword,
		Logger:  l.Logger.Copy(),
		stack:   l.copyStack(),
	}
}

func (l *traceLogger) Copy() Logger {
	return l.clone()
}

func (l *traceLogger) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey{}, l)
}

func (l *traceLogger) WithField(key string, value any) Logger {
	if key == l.keyword {
		stack := l.copyStack()
		stack[key] = append(stack[key], value)

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
		stackFields    = make(map[string]any)
		normalFields   = make(map[string]any, len(fields))
	)

	for k, v := range fields {
		if k == l.keyword {
			stackFields[k] = v
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
		stack = l.copyStack()
		for k, v := range stackFields {
			stack[k] = append(stack[k], v)
		}
	}

	return &traceLogger{
		keyword: l.keyword,
		Logger:  logger,
		stack:   stack,
	}
}

func (l *traceLogger) fieldsToAttach() map[string]any {
	l.stackMu.RLock()
	defer l.stackMu.RUnlock()

	if len(l.stack) == 0 {
		return nil
	}

	fields := make(map[string]any, len(l.stack))
	builder := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		stringBuilderPool.Put(builder)
	}()

	for k, v := range l.stack {
		if len(v) == 0 {
			continue
		}

		builder.Grow(len(v) * 10)

		for i, elem := range v {
			if i > 0 {
				builder.WriteString(" -> ")
			}

			if str, ok := elem.(string); ok {
				builder.WriteString(str)
			} else {
				builder.WriteString(fmt.Sprintf("%v", elem))
			}
		}

		fields[k] = builder.String()
		builder.Reset()
	}

	return fields
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
