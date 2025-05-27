package internal

import (
	"bytes"
)

const (
	colorBlack   = "30m"
	colorRed     = "31m"
	colorGreen   = "32m"
	colorYellow  = "33m"
	colorBlue    = "34m"
	colorMagenta = "35m"
	colorCyan    = "36m"
	colorWhite   = "37m"
)

const (
	colorReverseRed = "41m"
)

const (
	colorBrightBlack   = "90m"
	colorBrightRed     = "91m"
	colorBrightGreen   = "92m"
	colorBrightYellow  = "93m"
	colorBrightBlue    = "94m"
	colorBrightMagenta = "95m"
	colorBrightCyan    = "96m"
	colorBrightWhite   = "97m"
)

// ColorizeToBuffer 直接寫入 buffer，避免中間字串分配
func ColorizeToBuffer(buf *bytes.Buffer, s string, c string) {
	buf.WriteString(_colorPrefix)
	buf.WriteString(c)
	buf.WriteString(s)
	buf.WriteString(_colorReset)
}
