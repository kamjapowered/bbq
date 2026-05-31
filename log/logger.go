package log

import (
	"iter"

	f "github.com/kamjapowered/bbq/functional"
	"github.com/kamjapowered/bbq/queue"
)

// store log records and forward them to sinks
//
// records are kept in fifo order up to the configured maximum
//
// this type is not safe for concurrent use
//
//microwave:export
type Logger struct {
	records *queue.RingQueue[Record]
	sinks   []Sink
	dropped uint64
}

// create a logger that keeps at most maxRecords records
//
// if maxRecords is less than or equal to zero, records are not retained
//
// time: O(n) due to allocation
//
//microwave:export
func NewLogger(maxRecords int) *Logger {
	return &Logger{
		records: queue.NewRingQueue[Record](maxRecords),
		sinks:   make([]Sink, 0),
	}
}

// add a sink to receive future records
//
// nil sinks are ignored
//
// time: O(1) amortised
func (l *Logger) AddSink(s Sink) {
	if s != nil {
		l.sinks = append(l.sinks, s)
	}
}

// add a record to the logger
//
// when the record buffer is full, the oldest record is discarded
//
// time: O(1) amortised
func (l *Logger) Log(level LogLevel, msg string, opts ...f.Option[Record]) {
	r := Record{
		Level:   level,
		Message: msg,
		Fields:  make([]Field, 0),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&r)
		}
	}

	l.Push(r)

	for _, s := range l.sinks {
		s.Log(r)
	}
}

// add a prebuilt record to the logger without notifying sinks
//
// when the record buffer is full, the oldest record is discarded and the drop
// counter is incremented
//
// use to re-add or forward records; Log is the normal entry point
//
// time: O(1) amortised
func (l *Logger) Push(r Record) {
	if l.records.Full() {
		_, _ = l.records.Dequeue()
		l.dropped++
	}
	l.records.Enqueue(r)
}

// return the number of records currently buffered
//
// time: O(1)
func (l *Logger) Len() int {
	return l.records.Len()
}

// return the total number of records discarded because the buffer was full
//
// the counter is cumulative and never reset
//
// time: O(1)
func (l *Logger) Dropped() uint64 {
	return l.dropped
}

// iterate over stored records in fifo order
//
// iteration stops if yield returns false
//
// time: O(n)
func (l *Logger) Iter() iter.Seq[Record] {
	return l.records.Iter()
}

// iterate over stored records in fifo order while removing them
//
// iteration stops if yield returns false
//
// time: O(n)
func (l *Logger) Drain() iter.Seq[Record] {
	return l.records.Drain()
}

// removes all records while preserving capacity
//
// time 0(n)
func (l *Logger) Clear() {
	l.records.Clear()
}
