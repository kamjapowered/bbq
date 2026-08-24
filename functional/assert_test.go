package functional_test

import (
	"strings"
	"testing"

	f "kamjapowered.com/bbq/functional"
)

func recovered(t *testing.T, fn func()) string {
	t.Helper()

	var msg string
	func() {
		defer func() {
			if r := recover(); r != nil {
				msg, _ = r.(string)
			}
		}()
		fn()
	}()

	return msg
}

func TestAssertPasses(t *testing.T) {
	if msg := recovered(t, func() { f.Assert(true, "unreachable") }); msg != "" {
		t.Fatalf("a true condition panicked: %q", msg)
	}

	if msg := recovered(t, func() { f.AssertF(true, "unreachable %d", 1) }); msg != "" {
		t.Fatalf("a true condition panicked: %q", msg)
	}
}

func TestAssertFails(t *testing.T) {
	msg := recovered(t, func() { f.Assert(false, "boom ", 42) })

	if !strings.Contains(msg, "assertion error") {
		t.Fatalf("missing header: %q", msg)
	}
	if !strings.Contains(msg, "boom 42") {
		t.Fatalf("missing message: %q", msg)
	}
}

func TestAssertFFormats(t *testing.T) {
	msg := recovered(t, func() { f.AssertF(false, "want %d, got %d", 3, 4) })

	if !strings.Contains(msg, "want 3, got 4") {
		t.Fatalf("format not applied: %q", msg)
	}
}

// the location must be the asserting call site, not assert.go
func TestAssertReportsCallerLocation(t *testing.T) {
	msg := recovered(t, func() { f.Assert(false) })

	if !strings.Contains(msg, "assert_test.go") {
		t.Fatalf("location is not the caller: %q", msg)
	}
	if strings.Contains(msg, "functional/assert.go") {
		t.Fatalf("location points inside the assert package: %q", msg)
	}
	if !strings.Contains(msg, "no message") {
		t.Fatalf("empty message not defaulted: %q", msg)
	}
}
