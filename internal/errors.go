package internal

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/yanun0323/logs/internal/buffer"
	"github.com/yanun0323/logs/internal/colorize"
	"github.com/yanun0323/logs/internal/errors"
)

func extractErrors(err any) []slog.Attr {
	yanunErr, ok := err.(errors.Error)
	if !ok {
		return []slog.Attr{
			slog.String(KeyErr, fmt.Sprintf("%+v", err)),
		}
	}

	args := make([]slog.Attr, 0, len(yanunErr.Attributes())+len(yanunErr.Stack())+2)
	args = append(args, slog.String(KeyErr, yanunErr.Message()))
	args = append(args, slog.Any(KeyErrorsCause, yanunErr.Cause()))

	for _, a := range yanunErr.Attributes() {
		attr, ok := a.(errors.Attr)
		if !ok {
			println("not errors.Attr")
			continue
		}

		key, value := attr.Parameters()
		args = append(args, slog.String(key, fmt.Sprintf("%+v", value)))
	}

	buf := buffer.Pool.Get().(*bytes.Buffer)
	defer buffer.Put(buf)
	buf.Reset()
	buf.Grow(1024)

	for _, f := range yanunErr.Stack() {
		frame, ok := f.(errors.Frame)
		if !ok {
			println("not errors.Frame")
			continue
		}

		file, function, line := frame.Parameters()
		buf.WriteString("        ")
		colorize.Fprint(buf, colorize.ColorBrightBlue, "[", function, "]")
		buf.WriteByte(' ')
		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(line)
		buf.WriteByte('\n')
	}

	args = append(args, slog.String(KeyErrorsStack, buf.String()))

	return args
}
