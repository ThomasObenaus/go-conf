package interfaces

import "fmt"

type LogLevel string

const (
	LogLevel_Debug LogLevel = "Debug"
	LogLevel_Info  LogLevel = "Info"
	LogLevel_Warn  LogLevel = "Warn"
	LogLevel_Error LogLevel = "Error"
)

type LoggerFunc func(lvl LogLevel, formatString string, a ...interface{})

// DebugLogger logs all messages of level debug or above
var DebugLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf("[%s] %s", lvl, msg)
}

// InfoLogger logs all messages of level info or above
var InfoLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	if lvl == LogLevel_Debug {
		return
	}
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf("[%s] %s", lvl, msg)
}

// WarnLogger logs all messages of level warn or above
var WarnLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	if lvl == LogLevel_Debug || lvl == LogLevel_Info {
		return
	}
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf("[%s] %s", lvl, msg)
}

// ErrorLogger logs all messages of level error
var ErrorLogger = func(lvl LogLevel, formatString string, a ...interface{}) {
	if lvl != LogLevel_Error {
		return
	}
	msg := fmt.Sprintf(formatString, a...)
	fmt.Printf("[%s] %s", lvl, msg)
}

var NoLogging = func(lvl LogLevel, formatString string, a ...interface{}) {
}
