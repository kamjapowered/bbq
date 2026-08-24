package log

import (
	"errors"
	"slices"
	"testing"

	f "kamjapowered.com/bbq/functional"
)

// an in memory sink that keeps every record it receives
type fakeSink struct {
	records []Record
}

func (s *fakeSink) Log(r Record) {
	s.records = append(s.records, r)
}

func messages(t *testing.T, l *Logger) []string {
	t.Helper()

	out := make([]string, 0, l.Len())
	for r := range l.Iter() {
		out = append(out, r.Message)
	}

	return out
}

func TestNewLoggerStartsEmpty(t *testing.T) {
	l := NewLogger(8)

	if l.Len() != 0 {
		t.Errorf("Len = %d, want 0", l.Len())
	}
	if l.Dropped() != 0 {
		t.Errorf("Dropped = %d, want 0", l.Dropped())
	}
	if got := messages(t, l); len(got) != 0 {
		t.Errorf("Iter yielded %v, want nothing", got)
	}
}

func TestLogStoresRecord(t *testing.T) {
	l := NewLogger(4)

	l.Log(LogWarning, "cache miss", LogOpt.Field("key", "a"))

	if l.Len() != 1 {
		t.Fatalf("Len = %d, want 1", l.Len())
	}

	var got Record
	for r := range l.Iter() {
		got = r
	}

	if got.Level != LogWarning {
		t.Errorf("Level = %v, want %v", got.Level, LogWarning)
	}
	if got.Message != "cache miss" {
		t.Errorf("Message = %q, want %q", got.Message, "cache miss")
	}
	if len(got.Fields) != 1 || got.Fields[0].Key() != "key" {
		t.Errorf("Fields = %v, want one field keyed \"key\"", got.Fields)
	}
}

func TestLogWithoutOptionsHasEmptyFields(t *testing.T) {
	l := NewLogger(1)

	l.Log(LogInfo, "plain")

	for r := range l.Iter() {
		if r.Fields == nil {
			t.Fatal("Fields is nil, want an empty slice")
		}
		if len(r.Fields) != 0 {
			t.Fatalf("Fields = %v, want empty", r.Fields)
		}
	}
}

func TestLogIgnoresNilOptions(t *testing.T) {
	l := NewLogger(1)

	var nilOpt f.Option[Record]
	l.Log(LogInfo, "survives", nilOpt, LogOpt.Field("k", 1), nil)

	for r := range l.Iter() {
		if len(r.Fields) != 1 {
			t.Fatalf("Fields = %v, want one field", r.Fields)
		}
	}
}

func TestLogDispatchesToSinks(t *testing.T) {
	l := NewLogger(4)
	first, second := &fakeSink{}, &fakeSink{}

	l.AddSink(first)
	l.AddSink(second)

	l.Log(LogError, "boom", LogOpt.Wrap(errors.New("cause")))
	l.Log(LogDebug, "detail")

	for name, s := range map[string]*fakeSink{"first": first, "second": second} {
		if len(s.records) != 2 {
			t.Fatalf("%s sink got %d records, want 2", name, len(s.records))
		}
		if s.records[0].Message != "boom" || s.records[1].Message != "detail" {
			t.Errorf("%s sink got %q and %q, want \"boom\" and \"detail\"",
				name, s.records[0].Message, s.records[1].Message)
		}
		wrapped, _ := s.records[0].Fields[0].Value().(error)
		if len(s.records[0].Fields) != 1 || wrapped == nil || wrapped.Error() != "cause" {
			t.Errorf("%s sink lost the wrapped error: %v", name, s.records[0].Fields)
		}
	}
}

func TestSinksAddedLaterOnlySeeNewRecords(t *testing.T) {
	l := NewLogger(4)

	l.Log(LogInfo, "before")

	s := &fakeSink{}
	l.AddSink(s)

	l.Log(LogInfo, "after")

	if len(s.records) != 1 {
		t.Fatalf("sink got %d records, want 1", len(s.records))
	}
	if s.records[0].Message != "after" {
		t.Fatalf("sink got %q, want %q", s.records[0].Message, "after")
	}
}

func TestAddSinkIgnoresNil(t *testing.T) {
	l := NewLogger(2)
	s := &fakeSink{}

	l.AddSink(nil)
	l.AddSink(s)
	l.AddSink(nil)

	l.Log(LogInfo, "still delivered")

	if len(s.records) != 1 {
		t.Fatalf("sink got %d records, want 1", len(s.records))
	}
}

func TestPushDoesNotNotifySinks(t *testing.T) {
	l := NewLogger(2)
	s := &fakeSink{}
	l.AddSink(s)

	l.Push(Record{Level: LogInfo, Message: "forwarded"})

	if l.Len() != 1 {
		t.Fatalf("Len = %d, want 1", l.Len())
	}
	if len(s.records) != 0 {
		t.Fatalf("sink got %d records, want 0", len(s.records))
	}
}

func TestOverflowDropsOldest(t *testing.T) {
	l := NewLogger(2)

	for _, msg := range []string{"a", "b", "c", "d"} {
		l.Log(LogInfo, msg)
	}

	if l.Len() != 2 {
		t.Errorf("Len = %d, want 2", l.Len())
	}
	if l.Dropped() != 2 {
		t.Errorf("Dropped = %d, want 2", l.Dropped())
	}
	if got := messages(t, l); !slices.Equal(got, []string{"c", "d"}) {
		t.Errorf("records = %v, want [c d]", got)
	}
}

func TestOverflowStillReachesSinks(t *testing.T) {
	l := NewLogger(1)
	s := &fakeSink{}
	l.AddSink(s)

	l.Log(LogInfo, "a")
	l.Log(LogInfo, "b")

	if len(s.records) != 2 {
		t.Fatalf("sink got %d records, want 2", len(s.records))
	}
	if l.Len() != 1 || l.Dropped() != 1 {
		t.Fatalf("Len = %d and Dropped = %d, want 1 and 1", l.Len(), l.Dropped())
	}
}

func TestNonPositiveCapacityIsRaisedToOne(t *testing.T) {
	for _, capacity := range []int{0, -1} {
		l := NewLogger(capacity)
		s := &fakeSink{}
		l.AddSink(s)

		l.Log(LogInfo, "a")
		l.Log(LogInfo, "b")

		if l.Len() != 1 {
			t.Errorf("capacity %d: Len = %d, want 1", capacity, l.Len())
		}
		if got := messages(t, l); len(got) != 1 || got[0] != "b" {
			t.Errorf("capacity %d: Iter yielded %v, want [b]", capacity, got)
		}
		// one record fit, so only the first is dropped
		if l.Dropped() != 1 {
			t.Errorf("capacity %d: Dropped = %d, want 1", capacity, l.Dropped())
		}
		// sinks still see the records
		if len(s.records) != 2 {
			t.Errorf("capacity %d: sink got %d records, want 2", capacity, len(s.records))
		}
	}
}

func TestIterStopsEarly(t *testing.T) {
	l := NewLogger(4)
	for _, msg := range []string{"a", "b", "c"} {
		l.Log(LogInfo, msg)
	}

	seen := make([]string, 0, 2)
	for r := range l.Iter() {
		seen = append(seen, r.Message)
		if len(seen) == 2 {
			break
		}
	}

	if !slices.Equal(seen, []string{"a", "b"}) {
		t.Errorf("seen = %v, want [a b]", seen)
	}
	// breaking out does not consume records
	if l.Len() != 3 {
		t.Errorf("Len = %d, want 3", l.Len())
	}
}

func TestIterIsRepeatable(t *testing.T) {
	l := NewLogger(4)
	l.Log(LogInfo, "a")
	l.Log(LogInfo, "b")

	first := messages(t, l)
	second := messages(t, l)

	if !slices.Equal(first, second) {
		t.Fatalf("Iter yielded %v then %v", first, second)
	}
}

func TestDrainEmptiesTheLogger(t *testing.T) {
	l := NewLogger(4)
	for _, msg := range []string{"a", "b", "c"} {
		l.Log(LogInfo, msg)
	}

	drained := make([]string, 0, 3)
	for r := range l.Drain() {
		drained = append(drained, r.Message)
	}

	if !slices.Equal(drained, []string{"a", "b", "c"}) {
		t.Errorf("drained = %v, want [a b c]", drained)
	}
	if l.Len() != 0 {
		t.Errorf("Len = %d, want 0", l.Len())
	}
}

func TestDrainStopsEarlyAndKeepsTheRest(t *testing.T) {
	l := NewLogger(4)
	for _, msg := range []string{"a", "b", "c"} {
		l.Log(LogInfo, msg)
	}

	for range l.Drain() {
		break
	}

	if l.Len() != 2 {
		t.Fatalf("Len = %d, want 2", l.Len())
	}
	if got := messages(t, l); !slices.Equal(got, []string{"b", "c"}) {
		t.Fatalf("records = %v, want [b c]", got)
	}
}

func TestClearKeepsCapacityAndDropCount(t *testing.T) {
	l := NewLogger(2)
	for _, msg := range []string{"a", "b", "c"} {
		l.Log(LogInfo, msg)
	}

	l.Clear()

	if l.Len() != 0 {
		t.Errorf("Len = %d, want 0", l.Len())
	}
	// the drop counter is cumulative and survives a clear
	if l.Dropped() != 1 {
		t.Errorf("Dropped = %d, want 1", l.Dropped())
	}

	l.Log(LogInfo, "d")
	l.Log(LogInfo, "e")

	if got := messages(t, l); !slices.Equal(got, []string{"d", "e"}) {
		t.Errorf("records = %v, want [d e]", got)
	}
	if l.Dropped() != 1 {
		t.Errorf("Dropped = %d, want 1 after refilling to capacity", l.Dropped())
	}
}
