package csv

import (
	stdcsv "encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// the package only exposes path based helpers, so every case goes through a
// temporary file instead of an in memory reader
func write(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "data.csv")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing fixture: %v", err)
	}
	return path
}

func assertColumn(t *testing.T, data map[string][]any, header string, want []any) {
	t.Helper()
	got, ok := data[header]
	if !ok {
		t.Fatalf("missing column %q, got %v", header, data)
	}
	if len(got) != len(want) {
		t.Fatalf("column %q = %v, want %v", header, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("column %q[%d] = %v, want %v", header, i, got[i], want[i])
		}
	}
}

func TestReadWellFormed(t *testing.T) {
	path := write(t, "name,age,city\nalice,30,madrid\nbob,25,lisbon\n")

	data, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if len(data) != 3 {
		t.Fatalf("got %d columns, want 3", len(data))
	}
	assertColumn(t, data, "name", []any{"alice", "bob"})
	assertColumn(t, data, "age", []any{"30", "25"})
	assertColumn(t, data, "city", []any{"madrid", "lisbon"})
}

func TestReadKeepsEverythingAsStrings(t *testing.T) {
	path := write(t, "n,f,b\n1,2.5,true\n")

	data, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	// values are never converted, they stay strings inside the any slots
	for _, header := range []string{"n", "f", "b"} {
		if _, ok := data[header][0].(string); !ok {
			t.Errorf("column %q holds %T, want string", header, data[header][0])
		}
	}
}

func TestReadQuotedFields(t *testing.T) {
	path := write(t, "a,b\n\"x,y\",\"line1\nline2\"\n\"say \"\"hi\"\"\",plain\n")

	data, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	assertColumn(t, data, "a", []any{"x,y", `say "hi"`})
	assertColumn(t, data, "b", []any{"line1\nline2", "plain"})
}

func TestReadHeaderOnly(t *testing.T) {
	path := write(t, "a,b\n")

	data, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if len(data) != 2 {
		t.Fatalf("got %d columns, want 2", len(data))
	}
	for _, header := range []string{"a", "b"} {
		if len(data[header]) != 0 {
			t.Errorf("column %q = %v, want empty", header, data[header])
		}
	}
}

func TestReadDuplicateHeaders(t *testing.T) {
	path := write(t, "a,a\n1,2\n")

	data, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	// duplicated headers collapse into a single column holding both values
	if len(data) != 1 {
		t.Fatalf("got %d columns, want 1", len(data))
	}
	assertColumn(t, data, "a", []any{"1", "2"})
}

func TestReadEmptyFile(t *testing.T) {
	path := write(t, "")

	data, err := Read(path)
	if err == nil {
		t.Fatalf("Read of empty file returned %v, want an error", data)
	}
	if data != nil {
		t.Errorf("Read returned %v alongside the error, want nil", data)
	}
}

func TestReadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nope.csv")

	if _, err := Read(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Read of missing file = %v, want os.ErrNotExist", err)
	}
}

func TestReadRaggedRows(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"extra field", "a,b\n1,2,3\n"},
		{"missing field", "a,b\n1\n"},
		{"ragged later row", "a,b\n1,2\n3\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Read(write(t, tt.content))
			if err == nil {
				t.Fatalf("Read of ragged input returned %v, want an error", data)
			}
			var perr *stdcsv.ParseError
			if !errors.As(err, &perr) || !errors.Is(perr.Err, stdcsv.ErrFieldCount) {
				t.Errorf("Read = %v, want a field count parse error", err)
			}
		})
	}
}

func TestReadMalformedQuotes(t *testing.T) {
	path := write(t, "a,b\n\"unterminated,2\n")

	if _, err := Read(path); err == nil {
		t.Fatal("Read of malformed quotes returned no error")
	}
}

func TestSaveAndRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")
	data := map[string][]any{
		"name": {"alice", "bob"},
		"age":  {30, 25},
	}

	if err := Save(data, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	back, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	// everything comes back as strings, and column order is not preserved
	// because headers are collected from a map
	assertColumn(t, back, "name", []any{"alice", "bob"})
	assertColumn(t, back, "age", []any{"30", "25"})
}

func TestSaveFormatsValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")
	data := map[string][]any{
		"v": {"text", 7, 2.5, true, nil},
	}

	if err := Save(data, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	back, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	assertColumn(t, back, "v", []any{"text", "7", "2.500000", "true", "<nil>"})
}

func TestSaveQuotesAwkwardValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")
	data := map[string][]any{
		"v": {"a,b", "line1\nline2", `say "hi"`},
	}

	if err := Save(data, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	back, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	assertColumn(t, back, "v", []any{"a,b", "line1\nline2", `say "hi"`})
}

func TestSaveNoRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")
	data := map[string][]any{"a": {}, "b": {}}

	if err := Save(data, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if want := len("a,b\n"); len(content) != want {
		t.Errorf("file = %q, want just a header line", content)
	}
}

func TestSaveToUnwritablePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-dir", "out.csv")
	data := map[string][]any{"a": {1}}

	if err := Save(data, path); err == nil {
		t.Fatal("Save to a missing directory returned no error")
	}
}

// defect: Save indexes headers[0] before checking that there is a header, so
// an empty (or nil) map panics instead of writing an empty file or returning
// an error
func TestSaveEmptyMap(t *testing.T) {
	for _, tt := range []struct {
		name string
		data map[string][]any
	}{
		{"empty", map[string][]any{}},
		{"nil", nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "out.csv")

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Save(%s map) panicked: %v", tt.name, r)
				}
			}()

			if err := Save(tt.data, path); err != nil {
				t.Logf("Save returned %v", err)
			}
		})
	}
}

// defect: Save assumes every column is as long as the one that happens to be
// first in map iteration order, so unequal column lengths either panic or
// silently drop rows depending on a random order
func TestSaveRaggedColumns(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")
	data := map[string][]any{
		"a": {1, 2},
		"b": {3},
	}

	err := Save(data, path)
	if err == nil {
		t.Fatal("Save accepted unequal column lengths, want an error")
	}
	if !strings.Contains(err.Error(), "same length") {
		t.Errorf("error = %v, want it to explain the length mismatch", err)
	}

	// a rejected call must not leave a file behind
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Errorf("Save wrote %q despite rejecting the input", path)
	}
}
