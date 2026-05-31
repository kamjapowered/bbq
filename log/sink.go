package log

// receive log records from a logger
type Sink interface {
	Log(r Record)
}
