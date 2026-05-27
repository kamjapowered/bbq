package terminal

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// moves the cursor to the home position (row 1, column 1).
	MoveHome = "\x1b[H"

	// moves the cursor to the beginning of the current line (carriage return).
	MoveStart = "\r"

	// moves the cursor up by 1 row.
	MoveUp = "\x1b[A"

	// moves the cursor down by 1 row.
	MoveDown = "\x1b[B"

	// moves the cursor right by 1 column.
	MoveRight = "\x1b[C"

	// moves the cursor left by 1 column.
	MoveLeft = "\x1b[D"
)

const (
	// clears the screen from the cursor to the end of the screen.
	ClearScreenEnd = "\x1b[0J"

	// clears the screen from the cursor to the beginning of the screen.
	ClearScreenStart = "\x1b[1J"

	// clears the entire screen.
	ClearScreen = "\x1b[2J"

	// clears the line from the cursor to the end of the line.
	ClearLineEnd = "\x1b[0K"

	// clears the line from the cursor to the beginning of the line.
	ClearLineStart = "\x1b[1K"

	// clears the entire current line.
	ClearLine = "\x1b[2K"
)

const (
	// resets all text attributes (style, foreground, and background).
	StyleReset = "\x1b[0m"

	// enables bold or increased intensity.
	StyleBold = "\x1b[1m"

	// enables dim or decreased intensity.
	StyleDim = "\x1b[2m"

	// enables italic text.
	StyleItalic = "\x1b[3m"

	// enables underlined text.
	StyleUnderline = "\x1b[4m"

	// enables blinking text (support varies by terminal).
	StyleBlink = "\x1b[5m"

	// swaps foreground and background colors.
	StyleReverse = "\x1b[7m"

	// hides text (support varies by terminal).
	StyleHidden = "\x1b[8m"

	// enables strikethrough.
	StyleStrike = "\x1b[9m"
)

const (
	// sets the foreground color to black.
	FgBlack = "\x1b[30m"
	// sets the foreground color to red.
	FgRed = "\x1b[31m"
	// sets the foreground color to green.
	FgGreen = "\x1b[32m"
	// sets the foreground color to yellow.
	FgYellow = "\x1b[33m"
	// sets the foreground color to blue.
	FgBlue = "\x1b[34m"
	// sets the foreground color to magenta.
	FgMagenta = "\x1b[35m"
	// sets the foreground color to cyan.
	FgCyan = "\x1b[36m"
	// sets the foreground color to white.
	FgWhite = "\x1b[37m"

	// sets the foreground color to bright black (often gray).
	FgBlackBright = "\x1b[90m"
	// sets the foreground color to bright red.
	FgRedBright = "\x1b[91m"
	// sets the foreground color to bright green.
	FgGreenBright = "\x1b[92m"
	// sets the foreground color to bright yellow.
	FgYellowBright = "\x1b[93m"
	// sets the foreground color to bright blue.
	FgBlueBright = "\x1b[94m"
	// sets the foreground color to bright magenta.
	FgMagentaBright = "\x1b[95m"
	// sets the foreground color to bright cyan.
	FgCyanBright = "\x1b[96m"
	// sets the foreground color to bright white.
	FgWhiteBright = "\x1b[97m"
)

const (
	// sets the background color to black.
	BgBlack = "\x1b[40m"
	// sets the background color to red.
	BgRed = "\x1b[41m"
	// sets the background color to green.
	BgGreen = "\x1b[42m"
	// sets the background color to yellow.
	BgYellow = "\x1b[43m"
	// sets the background color to blue.
	BgBlue = "\x1b[44m"
	// sets the background color to magenta.
	BgMagenta = "\x1b[45m"
	// sets the background color to cyan.
	BgCyan = "\x1b[46m"
	// sets the background color to white.
	BgWhite = "\x1b[47m"

	// sets the background color to bright black.
	BgBlackBright = "\x1b[100m"
	// sets the background color to bright red.
	BgRedBright = "\x1b[101m"
	// sets the background color to bright green.
	BgGreenBright = "\x1b[102m"
	// sets the background color to bright yellow.
	BgYellowBright = "\x1b[103m"
	// sets the background color to bright blue.
	BgBlueBright = "\x1b[104m"
	// sets the background color to bright magenta.
	BgMagentaBright = "\x1b[105m"
	// sets the background color to bright cyan.
	BgCyanBright = "\x1b[106m"
	// sets the background color to bright white.
	BgWhiteBright = "\x1b[107m"
)

// returns a 24-bit (truecolor) foreground color escape sequence from a hex string.
//
// accepts strings in the form "rrggbb" or "#rrggbb" (case insensitive).
// if the input is not exactly 6 hex digits, Hex returns StyleReset.
func Hex(code string) string {
	code = strings.TrimPrefix(code, "#")
	if len(code) != 6 {
		return StyleReset
	}

	r, err1 := strconv.ParseUint(code[0:2], 16, 8)
	g, err2 := strconv.ParseUint(code[2:4], 16, 8)
	b, err3 := strconv.ParseUint(code[4:6], 16, 8)
	if err1 != nil || err2 != nil || err3 != nil {
		return StyleReset
	}

	return Rgb(int(r), int(g), int(b))
}

// returns a 24-bit (truecolor) foreground color escape sequence.
//
// values outside [0, 255] are clamped.
func Rgb(r, g, b int) string {
	r, g, b = clamp8(r), clamp8(g), clamp8(b)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// returns a 24-bit (truecolor) background color escape sequence from a hex string.
//
// accepts strings in the form "rrggbb" or "#rrggbb" (case insensitive).
// if the input is not exactly 6 hex digits, BgHex returns StyleReset.
func BgHex(code string) string {
	code = strings.TrimPrefix(code, "#")
	if len(code) != 6 {
		return StyleReset
	}

	r, err1 := strconv.ParseUint(code[0:2], 16, 8)
	g, err2 := strconv.ParseUint(code[2:4], 16, 8)
	b, err3 := strconv.ParseUint(code[4:6], 16, 8)
	if err1 != nil || err2 != nil || err3 != nil {
		return StyleReset
	}

	return BgRgb(int(r), int(g), int(b))
}

// returns a 24-bit (truecolor) background color escape sequence.
//
// values outside [0, 255] are clamped.
func BgRgb(r, g, b int) string {
	r, g, b = clamp8(r), clamp8(g), clamp8(b)
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

func clamp8(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
