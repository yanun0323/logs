package logs

import (
	"context"
	"sync/atomic"
	"time"
)

type tickerLogger struct {
	last                int64
	sender              atomic.Bool
	intervalMillisecond int64

	Logger
}

// NewTickerLogger creates a new ticker logger with the given interval and level.
//
// A ticker logger is a logger that logs messages only when the interval has passed,
// otherwise the messages will be dropped.
//
// If option is not provided, the logger will write to the os.Stdout with console format.
func NewTickerLogger(interval time.Duration, level Level, option ...*Option) Logger {
	itv := interval.Milliseconds()
	return &tickerLogger{
		last:                time.Now().UnixMilli() - itv,
		intervalMillisecond: itv,
		Logger:              New(level, option...),
	}
}

func (l *tickerLogger) canBeFire() bool {
	var (
		now  = time.Now().UnixMilli()
		last = atomic.LoadInt64(&l.last)
	)

	available := last + l.intervalMillisecond
	notReady := now <= available
	if notReady {
		return false
	}

	if l.sender.Swap(true) {
		return false
	}
	defer l.sender.Store(false)

	atomic.StoreInt64(&l.last, now)

	return true
}

func (l *tickerLogger) Copy() Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		Logger:              l.Logger.Copy(),
	}
}

func (l *tickerLogger) WithContext(ctx context.Context) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		Logger:              l.Logger.WithContext(ctx),
	}
}

func (l *tickerLogger) WithError(err error) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		Logger:              l.Logger.WithError(err),
	}
}

func (l *tickerLogger) WithField(key string, value any) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		Logger:              l.Logger.WithField(key, value),
	}
}

func (l *tickerLogger) WithFields(fields map[string]any) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		Logger:              l.Logger.WithFields(fields),
	}
}

func (l *tickerLogger) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey{}, l)
}

func (l *tickerLogger) Log(level Level, args ...any) {
	if l.canBeFire() {
		l.Logger.Log(level, args...)
	}
}

func (l *tickerLogger) Logf(level Level, format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Logf(level, format, args...)
	}
}

func (l *tickerLogger) Debug(args ...any) {
	if l.canBeFire() {
		l.Logger.Debug(args...)
	}
}

func (l *tickerLogger) Debugf(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Debugf(format, args...)
	}
}

func (l *tickerLogger) Info(args ...any) {
	if l.canBeFire() {
		l.Logger.Info(args...)
	}
}

func (l *tickerLogger) Infof(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Infof(format, args...)
	}
}

func (l *tickerLogger) Warn(args ...any) {
	if l.canBeFire() {
		l.Logger.Warn(args...)
	}
}

func (l *tickerLogger) Warnf(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Warnf(format, args...)
	}
}

func (l *tickerLogger) Error(args ...any) {
	if l.canBeFire() {
		l.Logger.Error(args...)
	}
}

func (l *tickerLogger) Errorf(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Errorf(format, args...)
	}
}

func (l *tickerLogger) Fatal(args ...any) {
	l.Logger.Fatal(args...)
}

func (l *tickerLogger) Fatalf(format string, args ...any) {
	l.Logger.Fatalf(format, args...)
}
