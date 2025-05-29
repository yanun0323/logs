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
	nextFireTime        int64

	Logger
}

// NewTickerLogger creates a new ticker logger with the given interval and level.
//
// A ticker logger is a logger that logs messages only when the interval has passed,
// otherwise the messages will be dropped.
//
// If option is not provided, the logger will write to the os.Stdout with console format.
func NewTickerLogger(level Level, interval time.Duration, option ...*Option) Logger {
	itv := interval.Milliseconds()
	now := time.Now().UnixMilli()
	return &tickerLogger{
		last:                now - itv,
		intervalMillisecond: itv,
		nextFireTime:        now,
		Logger:              New(level, option...),
	}
}

func (l *tickerLogger) canBeFire() bool {
	now := time.Now().UnixMilli()
	nextFire := atomic.LoadInt64(&l.nextFireTime)

	if now < nextFire {
		return false
	}

	if !l.sender.CompareAndSwap(false, true) {
		return false
	}

	now = time.Now().UnixMilli()
	nextFire = atomic.LoadInt64(&l.nextFireTime)

	if now < nextFire {
		l.sender.Store(false)
		return false
	}

	newNextFire := now + l.intervalMillisecond
	atomic.StoreInt64(&l.last, now)
	atomic.StoreInt64(&l.nextFireTime, newNextFire)

	l.sender.Store(false)

	return true
}

func (l *tickerLogger) Copy() Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		nextFireTime:        atomic.LoadInt64(&l.nextFireTime),
		Logger:              l.Logger.Copy(),
	}
}

func (l *tickerLogger) WithContext(ctx context.Context) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		nextFireTime:        atomic.LoadInt64(&l.nextFireTime),
		Logger:              l.Logger.WithContext(ctx),
	}
}

func (l *tickerLogger) WithError(err error) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		nextFireTime:        atomic.LoadInt64(&l.nextFireTime),
		Logger:              l.Logger.WithError(err),
	}
}

func (l *tickerLogger) WithField(key string, value any) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		nextFireTime:        atomic.LoadInt64(&l.nextFireTime),
		Logger:              l.Logger.WithField(key, value),
	}
}

func (l *tickerLogger) WithFields(args ...any) Logger {
	return &tickerLogger{
		last:                atomic.LoadInt64(&l.last),
		intervalMillisecond: l.intervalMillisecond,
		nextFireTime:        atomic.LoadInt64(&l.nextFireTime),
		Logger:              l.Logger.WithFields(args...),
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
