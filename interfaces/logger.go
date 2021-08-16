package interfaces

import "fmt"

// LogLevel is a type to specify the log level
type LogLevel string

const (
	// LogLevelDebug debug loglevel (contains messages of type debug, info, warn and error)
	LogLevelDebug LogLevel = "Debug"
	// LogLevelInfo info loglevel (contains messages of type info, warn and error)
	LogLevelInfo LogLevel = "Info"
	// LogLevelWarn waring loglevel (contains messages of type warn and error)
	LogLevelWarn LogLevel = "Warn"
	// LogLevelError error loglevel (contains messages of type error)
	LogLevelError LogLevel = "Error"
)

//nolint
const formatString = "[%s] %s"

// LoggerFunc is a function type that is called whenever the library wants to log a message. This type can be implemented to integrate your custom logging system.
type LoggerFunc func(lvl LogLevel, formatString string, a ...interface{})

// DebugLogger logs all messages of level debug or above
var DebugLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf(formatString, lvl, msg)
}

// InfoLogger logs all messages of level info or above
var InfoLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	if lvl == LogLevelDebug {
		return
	}
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf(formatString, lvl, msg)
}

// WarnLogger logs all messages of level warn or above
var WarnLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	if lvl == LogLevelDebug || lvl == LogLevelInfo {
		return
	}
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf(formatString, lvl, msg)
}

// ErrorLogger logs all messages of level error
var ErrorLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	if lvl != LogLevelError {
		return
	}
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf(formatString, lvl, msg)
}

// NoLogging logs no messages at all
var NoLogging = func(lvl LogLevel, formatString string, a ...interface{}) {
}
