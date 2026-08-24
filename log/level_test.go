package log

import "testing"

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level LogLevel
		want  string
	}{
		{LogDebug, "debug"},
		{LogInfo, "info"},
		{LogWarning, "warning"},
		{LogError, "error"},
		{LogFatal, "fatal"},
		{LogMaxLevel, "fatal"},
		{LogLevel(200), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("LogLevel(%d).String() = %q, want %q", uint8(tt.level), got, tt.want)
		}
	}
}

func TestLogLevelSeverityOrder(t *testing.T) {
	levels := []LogLevel{LogDebug, LogInfo, LogWarning, LogError, LogFatal}

	for i := 1; i < len(levels); i++ {
		if levels[i-1] >= levels[i] {
			t.Fatalf("%v is not below %v", levels[i-1], levels[i])
		}
	}

	if LogMaxLevel != LogFatal {
		t.Fatalf("LogMaxLevel = %v, want %v", LogMaxLevel, LogFatal)
	}

	// the zero value of the type is the lowest level
	var zero LogLevel
	if zero != LogDebug {
		t.Fatalf("zero LogLevel = %v, want %v", zero, LogDebug)
	}
}
