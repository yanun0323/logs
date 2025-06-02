package logs

import (
	"io"
	"log/slog"
	"os"

	"github.com/yanun0323/logs/internal"
)

// defaultOption is the default configuration used when no Option is provided.
// It uses console format and outputs to os.Stdout.
var defaultOption = &Option{
	Format: FormatConsole,
	Output: os.Stdout,
}

// Option represents the configuration options for creating loggers.
// It allows customization of the output format and destination.
type Option struct {
	// Format specifies the log output format.
	// Available formats: FormatConsole (default), FormatText, FormatJSON
	Format Format

	// Output specifies the destination writer for log output.
	// Defaults to os.Stdout if not specified.
	Output io.Writer
}

// createLoggerHandler creates an appropriate slog.Handler based on the Option configuration.
// It returns different handler types based on the specified Format:
// - FormatText: slog.NewTextHandler
// - FormatJSON: slog.NewJSONHandler
// - FormatConsole (default): custom handler of logs
func (opt *Option) createLoggerHandler(level Level) slog.Handler {
	switch opt.Format {
	case FormatText:
		return slog.NewTextHandler(opt.output(), &slog.HandlerOptions{
			Level: slog.Level(level),
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					switch a.Value.Kind() {
					case slog.KindTime:
						a.Value = slog.StringValue(a.Value.Time().Format(internal.GetDefaultTimeFormat()))
					}
				}

				return a
			},
		})
	case FormatJSON:
		return slog.NewJSONHandler(opt.output(), &slog.HandlerOptions{
			Level: slog.Level(level),
		})
	default:
		return internal.NewLoggerHandler(opt.output(), int8(level))
	}
}

func (opt *Option) output() io.Writer {
	if opt.Output == nil {
		return os.Stdout
	}
	return opt.Output
}
