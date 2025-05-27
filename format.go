package logs

// Format represents the output format type for logging.
type Format int8

const (
	// FormatConsole outputs logs in a colored, human-readable console format.
	// This is the default format with enhanced readability for development.
	FormatConsole Format = iota + 1

	// FormatText outputs logs in standard slog text format.
	// Format: key=value pairs separated by spaces.
	FormatText

	// FormatJSON outputs logs in JSON format.
	// Each log entry is a single JSON object on one line.
	FormatJSON
)
