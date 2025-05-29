package logs

import (
	"bytes"
	"context"

	"github.com/yanun0323/logs/internal"
)

type traceLogger struct {
	Logger

	keyword string
	stack   bytes.Buffer
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
	}
}

func (l *traceLogger) clone() *traceLogger {
	buf := l.stack.Bytes()
	stack := bytes.Buffer{}
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

func (l *traceLogger) WithField(key string, value any) Logger {
	if key == l.keyword {
		str := internal.ValueToString(value)
		stack := bytes.Buffer{}
		stack.Grow(l.stack.Len() + len(str) + len(_traceSep))
		stack.Write(l.stack.Bytes())
		stack.WriteString(_traceSep)
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
		stackValues    []any // 改用 slice 而非 map，減少分配
		normalFields   map[string]any
	)

	// 只在需要時分配 normalFields
	for k, v := range fields {
		if k == l.keyword {
			stackValues = append(stackValues, v)
			hasStackFields = true
		} else {
			if normalFields == nil {
				normalFields = make(map[string]any, len(fields))
			}
			normalFields[k] = v
		}
	}

	logger := l.Logger
	stack := l.stack

	if len(normalFields) != 0 {
		logger = logger.WithFields(normalFields)
	}

	if hasStackFields {
		currentLen := l.stack.Len()
		stack = bytes.Buffer{}

		// 更精確的容量預估：當前長度 + 分隔符數量 + 預估每個值的長度
		separatorCount := len(stackValues)
		if currentLen > 0 {
			separatorCount++ // 需要在當前內容後加分隔符
		} else {
			separatorCount-- // 第一個值前不需要分隔符
		}

		estimatedSize := currentLen + separatorCount*len(_traceSep) + len(stackValues)*16 // 16 是預估每個值的平均長度
		stack.Grow(estimatedSize)

		// 複製現有內容
		if currentLen > 0 {
			stack.Write(l.stack.Bytes())
		}

		for i, v := range stackValues {
			if currentLen > 0 || i > 0 {
				stack.WriteString(_traceSep)
			}
			stack.WriteString(internal.ValueToString(v))
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
	logger := l.Logger.WithFields(fields)

	return logger
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
	newLogger := l.Logger.WithError(err)
	if newLogger == l.Logger {
		// 如果底層 logger 沒有變化，直接返回自己
		return l
	}
	return &traceLogger{
		keyword: l.keyword,
		Logger:  newLogger,
		stack:   l.stack,
	}
}

func (l *traceLogger) WithContext(ctx context.Context) Logger {
	newLogger := l.Logger.WithContext(ctx)
	if newLogger == l.Logger {
		// 如果底層 logger 沒有變化，直接返回自己
		return l
	}
	return &traceLogger{
		keyword: l.keyword,
		Logger:  newLogger,
		stack:   l.stack,
	}
}

func (l *traceLogger) WithFunc(function string) Logger {
	newLogger := l.Logger.WithFunc(function)
	if newLogger == l.Logger {
		// 如果底層 logger 沒有變化，直接返回自己
		return l
	}
	return &traceLogger{
		keyword: l.keyword,
		Logger:  newLogger,
		stack:   l.stack,
	}
}
