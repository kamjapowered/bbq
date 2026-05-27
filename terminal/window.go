package terminal

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"golang.org/x/term"
)

// return whether the given file descriptor refers to a terminal.
func IsTerminal(fd int) bool {
	if fd < 0 {
		return false
	}
	return term.IsTerminal(fd)
}

// return the current terminal size in character cells for the given file descriptor.
//
// the first return value is the number of columns.
// the second return value is the number of rows.
// an error is returned if the descriptor is not a terminal or the size cannot be queried.
func GetSize(fd int) (int, int, error) {
	if fd < 0 {
		return 0, 0, errors.New("terminal: fd must be non negative")
	}

	cols, rows, err := term.GetSize(fd)
	if err != nil {
		return 0, 0, err
	}

	return cols, rows, nil
}

// return the current terminal size for standard input.
func GetSizeStdin() (int, int, error) {
	return GetSize(int(os.Stdin.Fd()))
}

// return the current terminal size for standard output.
func GetSizeStdout() (int, int, error) {
	return GetSize(int(os.Stdout.Fd()))
}

// return the current terminal size for standard error.
func GetSizeStderr() (int, int, error) {
	return GetSize(int(os.Stderr.Fd()))
}

// emit the terminal size each time the window is resized.
//
// on unix like systems, resize events are delivered via sigwinch.
// on windows, resize notifications are not available and an error is returned.
//
// the returned channel is closed when ctx is done.
// the first value sent is the current size if it can be determined.
func WatchResize(ctx context.Context) (<-chan [2]int, error) {
	if runtime.GOOS == "windows" {
		return nil, errors.New("terminal: resize notifications are not supported on windows")
	}

	out := make(chan [2]int, 1)

	go func() {
		defer close(out)

		// send initial size if available
		if cols, rows, err := GetSizeStdout(); err == nil && cols > 0 && rows > 0 {
			select {
			case out <- [2]int{cols, rows}:
			case <-ctx.Done():
				return
			}
		}

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		defer signal.Stop(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				cols, rows, err := GetSizeStdout()
				if err != nil || cols <= 0 || rows <= 0 {
					continue
				}

				// keep the most recent size if the receiver is slow
				select {
				case out <- [2]int{cols, rows}:
				default:
					select {
					case <-out:
					default:
					}
					select {
					case out <- [2]int{cols, rows}:
					default:
					}
				}
			}
		}
	}()

	return out, nil
}
