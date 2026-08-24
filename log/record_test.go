package log

import (
	"errors"
	"fmt"
	"testing"

	f "kamjapowered.com/bbq/functional"
)

// build a record the way Logger.Log does, without needing a logger
func build(opts ...f.Option[Record]) Record {
	r := Record{Fields: make([]Field, 0)}
	f.ApplyOptions(&r, opts...)
	return r
}

func TestFieldAccessors(t *testing.T) {
	r := build(LogOpt.Field("attempt", 3))

	if len(r.Fields) != 1 {
		t.Fatalf("Fields length = %d, want 1", len(r.Fields))
	}

	got := r.Fields[0]
	if got.Key() != "attempt" {
		t.Errorf("Key = %q, want %q", got.Key(), "attempt")
	}
	if got.Value() != 3 {
		t.Errorf("Value = %v, want 3", got.Value())
	}
	if got.String() != "attempt: 3" {
		t.Errorf("String = %q, want %q", got.String(), "attempt: 3")
	}
}

func TestFieldStringFormatting(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"string", "value", "k: value"},
		{"int", 42, "k: 42"},
		{"float", 1.5, "k: 1.5"},
		{"bool", true, "k: true"},
		{"nil", nil, "k: <nil>"},
		{"slice", []int{1, 2}, "k: [1 2]"},
		{"error", errors.New("boom"), "k: boom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := build(LogOpt.Field("k", tt.value))
			if got := r.Fields[0].String(); got != tt.want {
				t.Fatalf("String = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFieldZeroValue(t *testing.T) {
	var field Field

	if field.Key() != "" {
		t.Errorf("Key = %q, want empty", field.Key())
	}
	if field.Value() != nil {
		t.Errorf("Value = %v, want nil", field.Value())
	}
	if field.String() != ": <nil>" {
		t.Errorf("String = %q, want %q", field.String(), ": <nil>")
	}
}

func TestFieldF(t *testing.T) {
	r := build(LogOpt.FieldF("path", "%s:%d", "main.go", 12))

	if len(r.Fields) != 1 {
		t.Fatalf("Fields length = %d, want 1", len(r.Fields))
	}
	if got := r.Fields[0].Value(); got != "main.go:12" {
		t.Errorf("Value = %v, want %q", got, "main.go:12")
	}
	if got := r.Fields[0].String(); got != "path: main.go:12" {
		t.Errorf("String = %q, want %q", got, "path: main.go:12")
	}

	// a format with no operands is stored verbatim
	r = build(LogOpt.FieldF("k", "no operands"))
	if got := r.Fields[0].Value(); got != "no operands" {
		t.Errorf("Value = %v, want %q", got, "no operands")
	}
}

func TestWrap(t *testing.T) {
	err := errors.New("disk on fire")
	r := build(LogOpt.Wrap(err))

	if len(r.Fields) != 1 {
		t.Fatalf("Fields length = %d, want 1", len(r.Fields))
	}
	if r.Fields[0].Key() != "error" {
		t.Errorf("Key = %q, want %q", r.Fields[0].Key(), "error")
	}
	got, ok := r.Fields[0].Value().(error)
	if !ok {
		t.Fatalf("Value = %T, want error", r.Fields[0].Value())
	}
	if got != err {
		t.Errorf("Value = %v, want %v", got, err)
	}
}

func TestWrapNested(t *testing.T) {
	base := errors.New("base")
	wrapped := fmt.Errorf("layer: %w", base)

	r := build(LogOpt.Wrap(wrapped))

	// Wrap stores the error itself, so the chain survives for errors.Is
	got, ok := r.Fields[0].Value().(error)
	if !ok {
		t.Fatalf("Value = %T, want error", r.Fields[0].Value())
	}
	if !errors.Is(got, base) {
		t.Fatal("errors.Is could not reach the base error through the field")
	}
	if got.Error() != "layer: base" {
		t.Errorf("Error() = %q, want %q", got.Error(), "layer: base")
	}
}

func TestWrapNilIsIgnored(t *testing.T) {
	r := build(LogOpt.Wrap(nil))

	if len(r.Fields) != 0 {
		t.Fatalf("Fields length = %d, want 0", len(r.Fields))
	}
}

func TestOptionsApplyInOrder(t *testing.T) {
	r := build(
		LogOpt.Field("first", 1),
		LogOpt.Wrap(errors.New("second")),
		LogOpt.FieldF("third", "%v", 3),
	)

	want := []string{"first", "error", "third"}
	if len(r.Fields) != len(want) {
		t.Fatalf("Fields length = %d, want %d", len(r.Fields), len(want))
	}

	for i, key := range want {
		if r.Fields[i].Key() != key {
			t.Errorf("Fields[%d].Key = %q, want %q", i, r.Fields[i].Key(), key)
		}
	}
}

func TestOptionsAreReusable(t *testing.T) {
	opt := LogOpt.Field("shared", "value")

	a := build(opt)
	b := build(opt)

	if len(a.Fields) != 1 || len(b.Fields) != 1 {
		t.Fatalf("Fields lengths = %d and %d, want 1 and 1", len(a.Fields), len(b.Fields))
	}
	if a.Fields[0] != b.Fields[0] {
		t.Fatalf("reused option produced %v and %v", a.Fields[0], b.Fields[0])
	}
}
