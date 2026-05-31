package log

// the severity level used by the logger
//
//microwave:export
type LogLevel uint8

const (
	// intended for development diagnostics
	//microwave:export
	LogDebug LogLevel = iota

	// intended for high level engine or game state messages
	//microwave:export
	LogInfo

	// intended for unexpected but recoverable situations
	//microwave:export
	LogWarning

	// intended for failures that affect behaviour
	//microwave:export
	LogError

	// intended for unrecoverable failures. it does not exit the process, it only records the message
	//microwave:export
	LogFatal

	// the highest defined log level
	//microwave:export
	LogMaxLevel = LogFatal
)

// return the level name
func (l LogLevel) String() string {
	switch l {
	case LogDebug:
		return "debug"
	case LogInfo:
		return "info"
	case LogWarning:
		return "warning"
	case LogError:
		return "error"
	case LogFatal:
		return "fatal"
	default:
		return "unknown"
	}
}
