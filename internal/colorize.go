package internal

import (
	"bytes"
)

const (
	colorBlack   = "30"
	colorRed     = "31"
	colorGreen   = "32"
	colorYellow  = "33"
	colorBlue    = "34"
	colorMagenta = "35"
	colorCyan    = "36"
	colorWhite   = "37"
)

const (
	colorReverseRed = "41"
)

const (
	colorBrightBlack   = "90"
	colorBrightRed     = "91"
	colorBrightGreen   = "92"
	colorBrightYellow  = "93"
	colorBrightBlue    = "94"
	colorBrightMagenta = "95"
	colorBrightCyan    = "96"
	colorBrightWhite   = "97"
)

// ColorizeToBuffer 直接寫入 buffer，避免中間字串分配
func ColorizeToBuffer(buf *bytes.Buffer, s string, c string) {
	buf.WriteString("\x1b[")
	buf.WriteString(c)
	buf.WriteByte('m')
	buf.WriteString(s)
	buf.WriteString("\x1b[0m")
}
