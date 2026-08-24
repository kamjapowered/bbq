package terminal

import (
	"strconv"
	"testing"
)

func TestEscapeConstants(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		// cursor movement
		{"MoveHome", MoveHome, "\x1b[H"},
		{"MoveStart", MoveStart, "\r"},
		{"MoveUp", MoveUp, "\x1b[A"},
		{"MoveDown", MoveDown, "\x1b[B"},
		{"MoveRight", MoveRight, "\x1b[C"},
		{"MoveLeft", MoveLeft, "\x1b[D"},

		// clearing
		{"ClearScreenEnd", ClearScreenEnd, "\x1b[0J"},
		{"ClearScreenStart", ClearScreenStart, "\x1b[1J"},
		{"ClearScreen", ClearScreen, "\x1b[2J"},
		{"ClearLineEnd", ClearLineEnd, "\x1b[0K"},
		{"ClearLineStart", ClearLineStart, "\x1b[1K"},
		{"ClearLine", ClearLine, "\x1b[2K"},

		// styles
		{"StyleReset", StyleReset, "\x1b[0m"},
		{"StyleBold", StyleBold, "\x1b[1m"},
		{"StyleDim", StyleDim, "\x1b[2m"},
		{"StyleItalic", StyleItalic, "\x1b[3m"},
		{"StyleUnderline", StyleUnderline, "\x1b[4m"},
		{"StyleBlink", StyleBlink, "\x1b[5m"},
		{"StyleReverse", StyleReverse, "\x1b[7m"},
		{"StyleHidden", StyleHidden, "\x1b[8m"},
		{"StyleStrike", StyleStrike, "\x1b[9m"},

		// foreground
		{"FgBlack", FgBlack, "\x1b[30m"},
		{"FgRed", FgRed, "\x1b[31m"},
		{"FgGreen", FgGreen, "\x1b[32m"},
		{"FgYellow", FgYellow, "\x1b[33m"},
		{"FgBlue", FgBlue, "\x1b[34m"},
		{"FgMagenta", FgMagenta, "\x1b[35m"},
		{"FgCyan", FgCyan, "\x1b[36m"},
		{"FgWhite", FgWhite, "\x1b[37m"},
		{"FgBlackBright", FgBlackBright, "\x1b[90m"},
		{"FgRedBright", FgRedBright, "\x1b[91m"},
		{"FgGreenBright", FgGreenBright, "\x1b[92m"},
		{"FgYellowBright", FgYellowBright, "\x1b[93m"},
		{"FgBlueBright", FgBlueBright, "\x1b[94m"},
		{"FgMagentaBright", FgMagentaBright, "\x1b[95m"},
		{"FgCyanBright", FgCyanBright, "\x1b[96m"},
		{"FgWhiteBright", FgWhiteBright, "\x1b[97m"},

		// background
		{"BgBlack", BgBlack, "\x1b[40m"},
		{"BgRed", BgRed, "\x1b[41m"},
		{"BgGreen", BgGreen, "\x1b[42m"},
		{"BgYellow", BgYellow, "\x1b[43m"},
		{"BgBlue", BgBlue, "\x1b[44m"},
		{"BgMagenta", BgMagenta, "\x1b[45m"},
		{"BgCyan", BgCyan, "\x1b[46m"},
		{"BgWhite", BgWhite, "\x1b[47m"},
		{"BgBlackBright", BgBlackBright, "\x1b[100m"},
		{"BgRedBright", BgRedBright, "\x1b[101m"},
		{"BgGreenBright", BgGreenBright, "\x1b[102m"},
		{"BgYellowBright", BgYellowBright, "\x1b[103m"},
		{"BgBlueBright", BgBlueBright, "\x1b[104m"},
		{"BgMagentaBright", BgMagentaBright, "\x1b[105m"},
		{"BgCyanBright", BgCyanBright, "\x1b[106m"},
		{"BgWhiteBright", BgWhiteBright, "\x1b[107m"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestRgb(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b int
		want    string
	}{
		{"black", 0, 0, 0, "\x1b[38;2;0;0;0m"},
		{"white", 255, 255, 255, "\x1b[38;2;255;255;255m"},
		{"mixed", 12, 34, 56, "\x1b[38;2;12;34;56m"},
		{"clamps low", -1, -100, 0, "\x1b[38;2;0;0;0m"},
		{"clamps high", 256, 1000, 255, "\x1b[38;2;255;255;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rgb(tt.r, tt.g, tt.b); got != tt.want {
				t.Fatalf("Rgb(%d, %d, %d) = %q, want %q", tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestBgRgb(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b int
		want    string
	}{
		{"black", 0, 0, 0, "\x1b[48;2;0;0;0m"},
		{"white", 255, 255, 255, "\x1b[48;2;255;255;255m"},
		{"mixed", 12, 34, 56, "\x1b[48;2;12;34;56m"},
		{"clamps low", -5, 0, -255, "\x1b[48;2;0;0;0m"},
		{"clamps high", 300, 255, 999, "\x1b[48;2;255;255;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BgRgb(tt.r, tt.g, tt.b); got != tt.want {
				t.Fatalf("BgRgb(%d, %d, %d) = %q, want %q", tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestHex(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{"bare", "ff8800", "\x1b[38;2;255;136;0m"},
		{"hashed", "#ff8800", "\x1b[38;2;255;136;0m"},
		{"uppercase", "FF8800", "\x1b[38;2;255;136;0m"},
		{"mixed case", "Ff8800", "\x1b[38;2;255;136;0m"},
		{"black", "000000", "\x1b[38;2;0;0;0m"},
		{"empty", "", StyleReset},
		{"only hash", "#", StyleReset},
		{"too short", "fff", StyleReset},
		{"too long", "ff88000", StyleReset},
		{"double hash", "##ff8800", StyleReset},
		{"not hex", "gggggg", StyleReset},
		{"signed component", "+f8800a", StyleReset},
		{"spaced", " f8800a", StyleReset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hex(tt.code); got != tt.want {
				t.Fatalf("Hex(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestBgHex(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{"bare", "ff8800", "\x1b[48;2;255;136;0m"},
		{"hashed", "#ff8800", "\x1b[48;2;255;136;0m"},
		{"uppercase", "FF8800", "\x1b[48;2;255;136;0m"},
		{"white", "ffffff", "\x1b[48;2;255;255;255m"},
		{"empty", "", StyleReset},
		{"too short", "abcde", StyleReset},
		{"not hex", "zzzzzz", StyleReset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BgHex(tt.code); got != tt.want {
				t.Fatalf("BgHex(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestHexMatchesRgb(t *testing.T) {
	// the hex helpers are the string front end of the numeric ones
	for _, code := range []string{"000000", "0a141e", "ffffff", "7f7f7f"} {
		component := func(s string) int {
			v, err := strconv.ParseUint(s, 16, 8)
			if err != nil {
				t.Fatalf("parsing %q: %v", s, err)
			}
			return int(v)
		}

		r, g, b := component(code[0:2]), component(code[2:4]), component(code[4:6])

		if got, want := Hex(code), Rgb(r, g, b); got != want {
			t.Errorf("Hex(%q) = %q, want %q", code, got, want)
		}
		if got, want := BgHex(code), BgRgb(r, g, b); got != want {
			t.Errorf("BgHex(%q) = %q, want %q", code, got, want)
		}
	}
}
