package log

import (
	"sync/atomic"

	"github.com/open-horizon/edge-utilities/logger"
)

var singleton atomic.Pointer[logger.Logger]

// Init Initialize Logger
func Init(parameters logger.Parameters) error {
	log := &logger.Logger{Tracing: false}
	if err := log.Init(parameters); err != nil {
		return err
	}
	singleton.Store(log)
	return nil
}

// ParseDestinationsList parses a list of destinations
func ParseDestinationsList(destinations string) ([]bool, bool) {
	if log := singleton.Load(); log != nil {
		return log.ParseDestinationsList(destinations)
	}
	return nil, false
}

// Stop Logger
func Stop() {
	if log := singleton.Load(); log != nil {
		log.Stop()
	}
}

// IsLogging checks if the logging level if higher or equal to the level parameter
func IsLogging(level int) bool {
	if log := singleton.Load(); log != nil {
		return log.IsLogging(level)
	}
	return false
}

// Status log
func Status(format string, a ...interface{}) {
	if log := singleton.Load(); log != nil {
		log.Status(format, a...)
	}
}

// Fatal log
func Fatal(format string, a ...interface{}) {
	if log := singleton.Load(); log != nil {
		log.Fatal(format, a...)
	}
}

// Error log
func Error(format string, a ...interface{}) {
	if log := singleton.Load(); log != nil {
		log.Error(format, a...)
	}
}

// Warning log
func Warning(format string, a ...interface{}) {
	if log := singleton.Load(); log != nil {
		log.Warning(format, a...)
	}
}

// Info log
func Info(format string, a ...interface{}) {
	if log := singleton.Load(); log != nil {
		log.Info(format, a...)
	}
}

// Debug log
func Debug(format string, a ...interface{}) {
	if log := singleton.Load(); log != nil {
		log.Debug(format, a...)
	}
}

// Trace log
func Trace(format string, a ...interface{}) {
	if log := singleton.Load(); log != nil {
		log.Trace(format, a...)
	}
}

// Dump a struct to the log
func Dump(label string, a interface{}) {
	if log := singleton.Load(); log != nil {
		log.Dump(label, a)
	}
}
