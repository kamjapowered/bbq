package log

import (
	"fmt"
	f "github.com/kamjapowered/bbq/functional"
)

// ======
// record
// ======

// describe one log entry
//
//microwave:export
type Record struct {
	Level   LogLevel
	Message string
	Fields  []Field
}

// store structured metadata for a record
//
//microwave:export
type Field struct {
	key   string
	value any
}

// return the field formatted as key: value
func (f Field) String() string {
	return fmt.Sprintf("%s: %v", f.key, f.value)
}

// return the field key
//
// time: O(1)
func (f Field) Key() string {
	return f.key
}

// return the field value
//
// time: O(1)
func (f Field) Value() any {
	return f.value
}

// =======
// options
// =======
type logOpt struct{}

// add a structured field to a record
//
// time: O(1) amortised
func (logOpt) Field(key string, value any) f.Option[Record] {
	return func(r *Record) {
		r.Fields = append(r.Fields, Field{
			key:   key,
			value: value,
		})
	}
}

// add a formatted structured field to a record
//
// time: O(n) where n is the formatted output length
func (logOpt) FieldF(key string, format string, a ...any) f.Option[Record] {
	return func(r *Record) {
		r.Fields = append(r.Fields, Field{
			key:   key,
			value: fmt.Sprintf(format, a...),
		})
	}
}

// add an error field to a record
//
// nil errors are ignored
//
// time: O(1) amortised
func (logOpt) Wrap(err error) f.Option[Record] {
	return func(r *Record) {
		if err != nil {
			r.Fields = append(r.Fields, Field{
				key:   "error",
				value: err.Error(),
			})
		}
	}
}

// provide option helpers for log records
//
//microwave:export
var LogOpt logOpt
