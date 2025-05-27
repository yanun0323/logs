package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

const (
	KeyError   = "error"
	KeyContext = "context"
	KeyFunc    = "func"
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
		LevelInfo:  colorBrightGreen,
		LevelWarn:  colorBrightYellow,
		LevelError: colorBrightRed,
		LevelFatal: colorReverseRed,
	}
)

// 使用 buffer pool 來重用 buffer，減少記憶體分配
var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

type loggerHandler struct {
	isStd bool
	level *int8
	attrs []slog.Attr
	out   io.Writer
}

func NewLoggerHandler(w io.Writer, level int8) *loggerHandler {
	return &loggerHandler{
		level: &level,
		out:   w,
		attrs: make([]slog.Attr, 0),
		isStd: w == os.Stdout || w == os.Stderr,
	}
}

func (h *loggerHandler) Level() slog.Level {
	return slog.Level(*h.level)
}

func (h *loggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Level()
}

func (h *loggerHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	buf.Grow(256)

	timeStr := r.Time.Format(GetDefaultTimeFormat())
	if h.isStd {
		ColorizeToBuffer(buf, timeStr, colorBlack)
	} else {
		buf.WriteString(timeStr)
	}

	buf.WriteByte(' ')

	level := int8(r.Level)
	if h.isStd {
		ColorizeToBuffer(buf, LevelTitle(level), LevelColor(level))
	} else {
		buf.WriteString(LevelTitle(level))
	}

	buf.WriteByte(' ')

	// Message
	buf.WriteString(r.Message)
	buf.WriteString("  ")

	// 寫入預存的屬性
	h.writeAttrs(buf)

	buf.WriteByte('\n')

	_, err := h.out.Write(buf.Bytes())
	return err
}

func (h *loggerHandler) writeAttrs(buf *bytes.Buffer) {
	var str string
	for _, attr := range h.attrs {
		if h.isStd {
			if key, ok := fieldKeyCache[attr.Key]; ok {
				buf.WriteString(key)
			} else {
				buf.WriteString("\x1b[")
				buf.WriteString(colorMagenta)
				buf.WriteByte('m')
				buf.WriteByte('[')
				buf.WriteString(attr.Key)
				buf.WriteByte(']')
				buf.WriteString("\x1b[0m")
			}
		} else {
			buf.WriteByte('"')
			buf.WriteString(attr.Key)
			buf.WriteByte('"')
			buf.WriteByte('=')
		}

		buf.WriteByte(' ')
		if f, ok := attrValueFunc[attr.Value.Kind()]; ok {
			str = f(attr.Value)
		} else {
			str = fmt.Sprint(attr.Value.Any())
		}

		if h.isStd {
			ColorizeToBuffer(buf, str, colorBlack)
			buf.WriteString("  ")
		} else {
			buf.WriteString(str)
			buf.WriteByte(',')
		}
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
		isStd: h.isStd,
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
		KeyError:   "\x1b[" + colorRed + "m[" + KeyError + "]\x1b[0m",
		KeyContext: "\x1b[" + colorCyan + "m[" + KeyContext + "]\x1b[0m",
		KeyFunc:    "\x1b[" + colorBrightBlue + "m[" + KeyFunc + "]\x1b[0m",
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
