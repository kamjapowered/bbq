package errors

import (
	"fmt"
	"io"
	"strings"

	"github.com/vistormu/go-dsa/terminal"
)

type Kind string

type field struct {
	Key string
	Val any
}

type Error struct {
	Kind   Kind
	Msg    string
	fields []field
	err    error

	ansi bool
}

// create a new typed error
func New(kind Kind, msg string) *Error {
	return &Error{Kind: kind, Msg: msg, ansi: true}
}

// enable or disable ansi formatting for %+v printing
func (e *Error) Ansi(on bool) *Error {
	e.ansi = on
	return e
}

// attach one field
func (e *Error) With(key string, val any) *Error {
	e.fields = append(e.fields, field{Key: key, Val: val})
	return e
}

// wrap an underlying error
func (e *Error) Wrap(err error) *Error {
	if err == nil {
		return e
	}
	e.err = err
	return e
}

// return the wrapped error for errors.is and errors.as
func (e *Error) Unwrap() error {
	return e.err
}

// return a stable one line message
func (e *Error) Error() string {
	if e.Msg == "" {
		return string(e.Kind)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Msg)
}

// format controls how fmt prints the error
//
// %v   prints one line like Error()
//
// %+v  prints a multi line bullet format with the wrap chain
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.pretty())
			return
		}
		io.WriteString(s, e.Error())
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *Error) pretty() string {
	var b strings.Builder

	writeHeader(&b, e.Kind, e.Msg, e.ansi)

	for _, f := range e.fields {
		fmt.Fprintf(&b, "   |> %s: %v\n", f.Key, f.Val)
	}

	if e.err != nil {
		fmt.Fprintf(&b, "   |> caused by: %+v\n", e.err)
	}

	return b.String()
}

func writeHeader(b *strings.Builder, kind Kind, msg string, ansi bool) {
	if ansi {
		if msg != "" {
			fmt.Fprintf(b, "%s%s%s: %s\n", terminal.FgRed, kind, terminal.StyleReset, msg)
			return
		}
		fmt.Fprintf(b, "%s%s%s\n", terminal.FgRed, kind, terminal.StyleReset)
		return
	}

	if msg != "" {
		fmt.Fprintf(b, "%s: %s\n", kind, msg)
		return
	}
	fmt.Fprintf(b, "%s\n", kind)
}
