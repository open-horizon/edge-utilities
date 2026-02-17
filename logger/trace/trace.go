package trace

import (
	"sync/atomic"

	"github.com/open-horizon/edge-utilities/logger"
)

var singleton atomic.Pointer[logger.Logger]

// Init Initialize Logger
func Init(parameters logger.Parameters) error {
	trace := &logger.Logger{Tracing: true}
	if err := trace.Init(parameters); err != nil {
		return err
	}
	singleton.Store(trace)
	return nil
}

// ParseDestinationsList parses a list of destinations
func ParseDestinationsList(destinations string) ([]bool, bool) {
	if trace := singleton.Load(); trace != nil {
		return trace.ParseDestinationsList(destinations)
	}
	return nil, false
}

// Stop Logger
func Stop() {
	if trace := singleton.Load(); trace != nil {
		trace.Stop()
	}
}

// IsLogging checks if the logging level if higher or equal to the level parameter
func IsLogging(level int) bool {
	if trace := singleton.Load(); trace != nil {
		return trace.IsLogging(level)
	}
	return false
}

// Status log
func Status(format string, a ...interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Status(format, a...)
	}
}

// Fatal log
func Fatal(format string, a ...interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Fatal(format, a...)
	}
}

// Error log
func Error(format string, a ...interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Error(format, a...)
	}
}

// Warning log
func Warning(format string, a ...interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Warning(format, a...)
	}
}

// Info log
func Info(format string, a ...interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Info(format, a...)
	}
}

// Debug log
func Debug(format string, a ...interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Debug(format, a...)
	}
}

// Trace log
func Trace(format string, a ...interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Trace(format, a...)
	}
}

// Dump a struct to the log
func Dump(label string, a interface{}) {
	if trace := singleton.Load(); trace != nil {
		trace.Dump(label, a)
	}
}

// StackTrace will log the current stack trace
func StackTrace() {
	if trace := singleton.Load(); trace != nil {
		trace.StackTrace()
	}
}
