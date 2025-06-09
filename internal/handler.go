package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/yanun0323/logs/internal/buffer"
	"github.com/yanun0323/logs/internal/colorize"
)

const (
	KeyErr  = "error"
	KeyCtx  = "context"
	KeyFunc = "func"
)

const (
	_space        byte   = ' '
	_newline      byte   = '\n'
	_bracketOpen  string = "["
	_bracketClose string = "]"
	_colorPrefix  string = "\x1b["
	_colorReset   string = "\x1b[0m"
)

// 預計算的 level 字串，避免運行時計算
var (
	levelTitleCache = map[int8]string{
		LevelDebug: LevelDebugTitle,
		LevelInfo:  LevelInfoTitle,
		LevelWarn:  LevelWarnTitle,
		LevelError: LevelErrorTitle,
		LevelFatal: LevelFatalTitle,
	}

	levelColorCache = map[int8]colorize.Color{
		LevelDebug: colorize.ColorBlue,
		LevelInfo:  colorize.ColorGreen,
		LevelWarn:  colorize.ColorYellow,
		LevelError: colorize.ColorRed,
		LevelFatal: colorize.ColorRedReversed,
	}
)

type loggerHandler struct {
	level *int8
	attrs []slog.Attr
	out   io.Writer
}

func NewLoggerHandler(w io.Writer, level int8) slog.Handler {
	return &loggerHandler{
		level: &level,
		out:   w,
		attrs: make([]slog.Attr, 0),
	}
}

func (h *loggerHandler) Level() slog.Level {
	return slog.Level(*h.level)
}

func (h *loggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Level()
}

func (h *loggerHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := buffer.Pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer buffer.Pool.Put(buf)

	buf.Grow(256)

	timeStr := r.Time.Format(GetDefaultTimeFormat())
	colorize.Fprint(buf, colorize.ColorBlack, timeStr)
	buf.WriteByte(_space)

	level := int8(r.Level)
	colorize.Fprint(buf, LevelColor(level), LevelTitle(level))
	buf.WriteByte(_space)

	buf.WriteString(r.Message)
	buf.WriteByte(_space)
	buf.WriteByte(_space)

	h.writeAttrs(buf)

	buf.WriteByte(_newline)

	_, err := h.out.Write(buf.Bytes())
	return err
}

func (h *loggerHandler) writeAttrs(buf *bytes.Buffer) {
	var (
		str   string
		stack slog.Attr
	)

	for _, attr := range h.attrs {
		if attr.Key == KeyErrorsStack {
			stack = attr
			continue
		}

		writeFieldKey(buf, attr)

		buf.WriteByte(_space)
		if f, ok := attrValueFunc[attr.Value.Kind()]; ok {
			str = f(attr.Value)
		} else {
			str = fmt.Sprint(attr.Value.Any())
		}

		colorize.Fprint(buf, colorize.ColorBlack, str)
		buf.WriteByte(_space)
		buf.WriteByte(_space)
	}

	if len(stack.Key) != 0 {
		buf.WriteByte('\n')
		buf.WriteString("    ")
		writeFieldKey(buf, stack)
		buf.WriteByte('\n')
		buf.WriteString(stack.Value.String())
	}
}

func writeFieldKey(buf *bytes.Buffer, attr slog.Attr) {
	if key, ok := fieldKeyCache[attr.Key]; ok {
		buf.WriteString(key)
	} else {
		colorize.Fprint(buf, colorize.ColorMagenta, _bracketOpen, attr.Key, _bracketClose)
	}
}

var attrValueFunc = map[slog.Kind]func(slog.Value) string{
	slog.KindString: func(v slog.Value) string {
		return v.String()
	},
	slog.KindInt64: func(v slog.Value) string {
		return strconv.FormatInt(v.Int64(), 10)
	},
	slog.KindUint64: func(v slog.Value) string {
		return strconv.FormatUint(v.Uint64(), 10)
	},
	slog.KindFloat64: func(v slog.Value) string {
		return strconv.FormatFloat(v.Float64(), 'g', -1, 64)
	},
	slog.KindBool: func(v slog.Value) string {
		return strconv.FormatBool(v.Bool())
	},
	slog.KindAny: func(v slog.Value) string {
		return fmt.Sprintf("%+v", v.Any())
	},
}

func (h *loggerHandler) clone() *loggerHandler {
	newAttrs := make([]slog.Attr, len(h.attrs))
	copy(newAttrs, h.attrs)
	return &loggerHandler{
		level: h.level,
		out:   h.out,
		attrs: newAttrs,
	}
}

func (h *loggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	hh := h.clone()
	hh.attrs = append(hh.attrs, attrs...)

	return hh
}

func (h *loggerHandler) WithGroup(name string) slog.Handler {
	return h
}

const (
	KeyErrorsStack = "error.stack"
	KeyErrorsCause = "error.cause"
)

var (
	fieldKeyCache = map[string]string{
		KeyErr:         colorize.Sprint(colorize.ColorRed, _bracketOpen, KeyErr, _bracketClose),
		KeyCtx:         colorize.Sprint(colorize.ColorCyan, _bracketOpen, KeyCtx, _bracketClose),
		KeyFunc:        colorize.Sprint(colorize.ColorBrightBlue, _bracketOpen, KeyFunc, _bracketClose),
		KeyErrorsCause: colorize.Sprint(colorize.ColorYellow, _bracketOpen, KeyErrorsCause, _bracketClose),
		KeyErrorsStack: colorize.Sprint(colorize.ColorCyan, _bracketOpen, KeyErrorsStack, _bracketClose),
	}
)

const (
	LevelFatal int8 = 12
	LevelError int8 = 8
	LevelWarn  int8 = 4
	LevelInfo  int8 = 0
	LevelDebug int8 = -4
)

const (
	LevelFatalTitle = "FATAL"
	LevelErrorTitle = "ERROR"
	LevelWarnTitle  = "WARN "
	LevelInfoTitle  = "INFO "
	LevelDebugTitle = "DEBUG"
)

func LevelTitle(level int8) string {
	if title, ok := levelTitleCache[level]; ok {
		return title
	}
	return LevelInfoTitle
}

func LevelColor(level int8) colorize.Color {
	if color, ok := levelColorCache[level]; ok {
		return color
	}
	return colorize.ColorBrightGreen
}
