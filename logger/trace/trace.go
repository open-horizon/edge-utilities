package trace

import (
	"github.com/open-horizon/edge-utilities/logger"
)

var trace = logger.Logger{Tracing: true, Logger: nil, Level: 0}

// Init Initialize Logger
func Init(parameters logger.Parameters) error {
	return trace.Init(parameters)
}

// Stop Logger
func Stop() {
	trace.Stop()
}

// IsLogging checks if the logging level if higher or equal to the level parameter
func IsLogging(level int) bool {
	return trace.IsLogging(level)
}

// Status log
func Status(format string, a ...interface{}) {
	trace.Status(format, a...)
}

// Fatal log
func Fatal(format string, a ...interface{}) {
	trace.Fatal(format, a...)
}

// Error log
func Error(format string, a ...interface{}) {
	trace.Error(format, a...)
}

// Warning log
func Warning(format string, a ...interface{}) {
	trace.Warning(format, a...)
}

// Info log
func Info(format string, a ...interface{}) {
	trace.Info(format, a...)
}

// Debug log
func Debug(format string, a ...interface{}) {
	trace.Debug(format, a...)
}

// Trace log
func Trace(format string, a ...interface{}) {
	trace.Trace(format, a...)
}

// Dump a struct to the log
func Dump(label string, a interface{}) {
	trace.Dump(label, a)
}

// StackTrace will log the current stack trace
func StackTrace() {
	trace.StackTrace()
}
