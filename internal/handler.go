package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/yanun0323/logs/internal/buffer"
)

const (
	KeyError   = "error"
	KeyContext = "context"
	KeyFunc    = "func"
)

const (
	_space        byte   = ' '
	_newline      byte   = '\n'
	_bracketOpen  byte   = '['
	_bracketClose byte   = ']'
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

	levelColorCache = map[int8]string{
		LevelDebug: colorBrightBlue,
		LevelInfo:  colorGreen,
		LevelWarn:  colorBrightYellow,
		LevelError: colorBrightRed,
		LevelFatal: colorReverseRed,
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
	ColorizeToBuffer(buf, timeStr, colorBlack)
	buf.WriteByte(_space)

	level := int8(r.Level)
	ColorizeToBuffer(buf, LevelTitle(level), LevelColor(level))
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
		str string
	)
	for _, attr := range h.attrs {
		if key, ok := fieldKeyCache[attr.Key]; ok {
			buf.WriteString(key)
		} else {
			buf.WriteString(_colorPrefix)
			buf.WriteString(colorMagenta)
			buf.WriteByte(_bracketOpen)
			buf.WriteString(attr.Key)
			buf.WriteByte(_bracketClose)
			buf.WriteString(_colorReset)
		}

		buf.WriteByte(_space)
		if f, ok := attrValueFunc[attr.Value.Kind()]; ok {
			str = f(attr.Value)
		} else {
			str = fmt.Sprint(attr.Value.Any())
		}

		ColorizeToBuffer(buf, str, colorBlack)
		buf.WriteByte(_space)
		buf.WriteByte(_space)
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
		return fmt.Sprint(v.Any())
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

var (
	fieldKeyCache = map[string]string{
		KeyError:   _colorPrefix + colorRed + "[" + KeyError + "]" + _colorReset,
		KeyContext: _colorPrefix + colorCyan + "[" + KeyContext + "]" + _colorReset,
		KeyFunc:    _colorPrefix + colorBrightBlue + "[" + KeyFunc + "]" + _colorReset,
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

func LevelColor(level int8) string {
	if color, ok := levelColorCache[level]; ok {
		return color
	}
	return colorBrightGreen
}
