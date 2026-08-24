package functional

import (
	"fmt"
	"runtime"
)

// panic with a located message when cond is false
//
// the operands are formatted like fmt.Sprint. the reported location is the
// caller's, not this package's, so the panic points at the assertion itself
//
// time: O(n) where n is the formatted message length
//
//microwave:export
func Assert(cond bool, a ...any) {
	if cond {
		return
	}

	panic(assertion(fmt.Sprint(a...)))
}

// panic with a located, formatted message when cond is false
//
// the format string version of Assert
//
// time: O(n) where n is the formatted message length
//
//microwave:export
func AssertF(cond bool, format string, a ...any) {
	if cond {
		return
	}

	panic(assertion(fmt.Sprintf(format, a...)))
}

func assertion(msg string) string {
	// 0 is this function, 1 is Assert/AssertF, 2 is the caller that asserted
	location := "unknown"
	if _, file, line, ok := runtime.Caller(2); ok {
		location = fmt.Sprintf("%s:%d", file, line)
	}

	if msg == "" {
		msg = "no message"
	}

	return fmt.Sprintf(
		"-> assertion error\n   |> location: %s\n   |> message: %s",
		location, msg,
	)
}
