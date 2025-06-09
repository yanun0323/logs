package logs

import (
	"bytes"
	"fmt"

	"github.com/yanun0323/logs/internal"
	"github.com/yanun0323/logs/internal/buffer"
	"github.com/yanun0323/logs/internal/colorize"
	"github.com/yanun0323/logs/internal/errors"
)

func extractErrors(err error) []any {
	yanunErr, ok := err.(errors.Error)
	if !ok {
		return []any{KeyErr, fmt.Sprintf("%+v", err)}
	}

	args := make([]any, 0, len(yanunErr.Attributes())+len(yanunErr.Stack())+2)
	args = append(args, internal.KeyErr, yanunErr.Message())
	args = append(args, internal.KeyErrorsCause, yanunErr.Cause())

	for _, a := range yanunErr.Attributes() {
		attr, ok := a.(errors.Attr)
		if !ok {
			println("not errors.Attr")
			continue
		}

		key, value := attr.Parameters()
		args = append(args, key, value)
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

	args = append(args, internal.KeyErrorsStack, buf.String())

	return args
}
