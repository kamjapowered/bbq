package terminal

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// a regular file descriptor, which is never a terminal
func regularFd(t *testing.T) int {
	t.Helper()

	path := filepath.Join(t.TempDir(), "not-a-terminal")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("creating the temp file: %v", err)
	}
	t.Cleanup(func() { file.Close() })

	return int(file.Fd())
}

func TestIsTerminalRejectsNegativeFd(t *testing.T) {
	for _, fd := range []int{-1, -42} {
		if IsTerminal(fd) {
			t.Errorf("IsTerminal(%d) = true, want false", fd)
		}
	}
}

func TestIsTerminalOnRegularFile(t *testing.T) {
	if IsTerminal(regularFd(t)) {
		t.Fatal("IsTerminal on a regular file = true, want false")
	}
}

func TestSizeRejectsNegativeFd(t *testing.T) {
	cols, rows, err := Size(-1)

	if err == nil {
		t.Fatal("Size(-1) returned no error")
	}
	if !strings.Contains(err.Error(), "terminal:") {
		t.Errorf("error = %q, want it prefixed by the package name", err.Error())
	}
	if cols != 0 || rows != 0 {
		t.Errorf("Size(-1) = %d, %d, want 0, 0", cols, rows)
	}
}

func TestSizeOnNonTerminalFails(t *testing.T) {
	// a regular file has no window size, so the ioctl must fail
	cols, rows, err := Size(regularFd(t))

	if err == nil {
		t.Fatal("Size on a regular file returned no error")
	}
	if cols != 0 || rows != 0 {
		t.Errorf("Size on a regular file = %d, %d, want 0, 0", cols, rows)
	}
}

// the std helpers must report the same thing as Size on the matching descriptor.
// the values themselves depend on how the test binary was started, so only the
// agreement is asserted
func TestSizeStdHelpersDelegate(t *testing.T) {
	tests := []struct {
		name string
		file *os.File
		fn   func() (int, int, error)
	}{
		{"SizeStdin", os.Stdin, SizeStdin},
		{"SizeStdout", os.Stdout, SizeStdout},
		{"SizeStderr", os.Stderr, SizeStderr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantCols, wantRows, wantErr := Size(int(tt.file.Fd()))
			cols, rows, err := tt.fn()

			if (err == nil) != (wantErr == nil) {
				t.Fatalf("%s error = %v, want %v", tt.name, err, wantErr)
			}
			if err == nil && (cols != wantCols || rows != wantRows) {
				t.Fatalf("%s = %d, %d, want %d, %d", tt.name, cols, rows, wantCols, wantRows)
			}
		})
	}
}

func TestWatchResizeClosesOnCancel(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("resize notifications are not supported on windows")
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch, err := WatchResize(ctx)
	if err != nil {
		t.Fatalf("WatchResize returned %v", err)
	}
	if ch == nil {
		t.Fatal("WatchResize returned a nil channel")
	}

	cancel()

	deadline := time.After(5 * time.Second)
	for {
		select {
		case size, ok := <-ch:
			if !ok {
				return
			}
			// the initial size may arrive before the cancellation is seen
			if size[0] <= 0 || size[1] <= 0 {
				t.Fatalf("received size %v, want positive values", size)
			}
		case <-deadline:
			t.Fatal("the channel was not closed after the context was cancelled")
		}
	}
}

func TestWatchResizeWithCancelledContext(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("resize notifications are not supported on windows")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch, err := WatchResize(ctx)
	if err != nil {
		t.Fatalf("WatchResize returned %v", err)
	}

	deadline := time.After(5 * time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("the channel was not closed for an already cancelled context")
		}
	}
}
