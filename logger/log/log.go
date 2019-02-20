package log

import (
	"github.com/open-horizon/edge-utilities/logger"
)

var log = logger.Logger{}

// Init Initialize Logger
func Init(parameters logger.Parameters) error {
	return log.Init(parameters)
}

// Stop Logger
func Stop() {
	log.Stop()
}

// IsLogging checks if the logging level if higher or equal to the level parameter
func IsLogging(level int) bool {
	return log.IsLogging(level)
}

// Status log
func Status(format string, a ...interface{}) {
	log.Status(format, a...)
}

// Fatal log
func Fatal(format string, a ...interface{}) {
	log.Fatal(format, a...)
}

// Error log
func Error(format string, a ...interface{}) {
	log.Error(format, a...)
}

// Warning log
func Warning(format string, a ...interface{}) {
	log.Warning(format, a...)
}

// Info log
func Info(format string, a ...interface{}) {
	log.Info(format, a...)
}

// Debug log
func Debug(format string, a ...interface{}) {
	log.Debug(format, a...)
}

// Trace log
func Trace(format string, a ...interface{}) {
	log.Trace(format, a...)
}

// Dump a struct to the log
func Dump(label string, a interface{}) {
	log.Dump(label, a)
}
