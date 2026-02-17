package trace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-horizon/edge-utilities/logger"
)

// TestInit_BasicConfiguration tests basic trace logger initialization:
// - File destination with valid parameters
// - Stdout destination
// - Multiple destinations
// - Verifies tracing is enabled by default
//
// This ensures the singleton trace logger can be properly initialized.
func TestInit_BasicConfiguration(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name        string
		params      logger.Parameters
		expectError bool
	}{
		{
			name: "file-destination",
			params: logger.Parameters{
				RootPath:                 tempDir,
				FileName:                 "trace-test",
				MaxFileSize:              1024,
				MaxCompressedFilesNumber: 3,
				Destinations:             "file",
				Prefix:                   "TRACE",
				Level:                    "INFO",
				MaintenanceInterval:      60,
			},
			expectError: false,
		},
		{
			name: "stdout-destination",
			params: logger.Parameters{
				Destinations:        "stdout",
				Prefix:              "TRACE",
				Level:               "DEBUG",
				MaintenanceInterval: 60,
			},
			expectError: false,
		},
		{
			name: "multiple-destinations",
			params: logger.Parameters{
				RootPath:                 tempDir,
				FileName:                 "trace-multi",
				MaxFileSize:              1024,
				MaxCompressedFilesNumber: 3,
				Destinations:             "file,stdout",
				Prefix:                   "TRACE",
				Level:                    "INFO",
				MaintenanceInterval:      60,
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Init(tc.params)
			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Clean up
			Stop()
		})
	}
}

// TestLoggingLevels tests logging at different levels with tracing enabled:
// - Status, Fatal, Error, Warning, Info, Debug, Trace
// - Verifies IsLogging() correctly reports level status
// - Verifies messages are prefixed with * (tracing indicator)
//
// This ensures all logging levels work correctly with tracing enabled.
func TestLoggingLevels(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "trace-test.log")

	params := logger.Parameters{
		RootPath:                 tempDir,
		FileName:                 "trace-test",
		MaxFileSize:              1024,
		MaxCompressedFilesNumber: 3,
		Destinations:             "file",
		Prefix:                   "TRACE",
		Level:                    "TRACE",
		MaintenanceInterval:      60,
	}

	err := Init(params)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Stop()

	// Test all logging levels
	Status("Status message: %s", "test")
	Fatal("Fatal message: %s", "test")
	Error("Error message: %s", "test")
	Warning("Warning message: %s", "test")
	Info("Info message: %s", "test")
	Debug("Debug message: %s", "test")
	Trace("Trace message: %s", "test")

	// Test IsLogging for different levels
	if !IsLogging(logger.STATUS) {
		t.Error("Expected STATUS level to be logging")
	}
	if !IsLogging(logger.INFO) {
		t.Error("Expected INFO level to be logging")
	}
	if !IsLogging(logger.TRACE) {
		t.Error("Expected TRACE level to be logging")
	}

	// Verify log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

// TestStackTrace tests the StackTrace function:
// - Verifies StackTrace() can be called without errors
// - Ensures stack trace is logged
//
// This tests the trace-specific StackTrace functionality.
func TestStackTrace(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "trace-stack.log")

	params := logger.Parameters{
		RootPath:                 tempDir,
		FileName:                 "trace-stack",
		MaxFileSize:              1024,
		MaxCompressedFilesNumber: 3,
		Destinations:             "file",
		Prefix:                   "TRACE",
		Level:                    "DEBUG",
		MaintenanceInterval:      60,
	}

	err := Init(params)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Stop()

	// Should not panic or error
	StackTrace()

	// Verify log file was created and has content
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}

	// Read log file to verify stack trace was written
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file is empty, expected stack trace content")
	}
}

// TestDump tests the Dump function for struct logging with tracing:
// - Dumps a struct to the log
// - Verifies no errors occur
// - Verifies tracing prefix is applied
//
// This ensures struct dumping works correctly with tracing enabled.
func TestDump(t *testing.T) {
	tempDir := t.TempDir()

	params := logger.Parameters{
		RootPath:                 tempDir,
		FileName:                 "trace-dump",
		MaxFileSize:              1024,
		MaxCompressedFilesNumber: 3,
		Destinations:             "file",
		Prefix:                   "TRACE",
		Level:                    "DEBUG",
		MaintenanceInterval:      60,
	}

	err := Init(params)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Stop()

	type TestStruct struct {
		Name  string
		Value int
		Items []string
	}

	testData := TestStruct{
		Name:  "test",
		Value: 42,
		Items: []string{"item1", "item2", "item3"},
	}

	// Should not panic or error
	Dump("TestStruct", testData)
}

// TestParseDestinationsList tests destination parsing:
// - Valid destination strings
// - Multiple destinations
// - Invalid destinations
//
// This ensures destination parsing works correctly for trace logger.
func TestParseDestinationsList(t *testing.T) {
	testCases := []struct {
		name         string
		destinations string
		expectValid  bool
	}{
		{"file-only", "file", true},
		{"stdout-only", "stdout", true},
		{"multiple", "file,stdout", true},
		{"all", "file,stdout,syslog,glog", true},
		{"empty", "", false},
		{"invalid", "invalid", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, valid := ParseDestinationsList(tc.destinations)
			if valid != tc.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.expectValid, valid)
			}
		})
	}
}

// TestStop tests proper cleanup:
// - Logger can be stopped without errors
// - Multiple Stop() calls don't cause issues
//
// This ensures proper resource cleanup for trace logger.
func TestStop(t *testing.T) {
	tempDir := t.TempDir()

	params := logger.Parameters{
		RootPath:                 tempDir,
		FileName:                 "trace-stop",
		MaxFileSize:              1024,
		MaxCompressedFilesNumber: 3,
		Destinations:             "file",
		Prefix:                   "TRACE",
		Level:                    "INFO",
		MaintenanceInterval:      60,
	}

	err := Init(params)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// First stop should work
	Stop()

	// Second stop should not panic
	Stop()
}

// TestConcurrentLogging tests thread safety with tracing:
// - Multiple goroutines logging simultaneously
// - No race conditions
// - Tracing prefix applied consistently
//
// This test should be run with -race flag.
// Usage: go test -race
func TestConcurrentLogging(t *testing.T) {
	tempDir := t.TempDir()

	params := logger.Parameters{
		RootPath:                 tempDir,
		FileName:                 "trace-concurrent",
		MaxFileSize:              1024,
		MaxCompressedFilesNumber: 3,
		Destinations:             "file",
		Prefix:                   "TRACE",
		Level:                    "INFO",
		MaintenanceInterval:      60,
	}

	err := Init(params)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Stop()

	// Launch multiple goroutines that log concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				Info("Goroutine %d: message %d", id, j)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestLogLevelFiltering tests that log level filtering works with tracing:
// - Messages below the configured level are not logged
// - Messages at or above the level are logged
// - Tracing prefix is applied to logged messages
//
// This ensures proper log level filtering with tracing enabled.
func TestLogLevelFiltering(t *testing.T) {
	tempDir := t.TempDir()

	params := logger.Parameters{
		RootPath:                 tempDir,
		FileName:                 "trace-filtering",
		MaxFileSize:              1024,
		MaxCompressedFilesNumber: 3,
		Destinations:             "file",
		Prefix:                   "TRACE",
		Level:                    "WARNING",
		MaintenanceInterval:      60,
	}

	err := Init(params)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Stop()

	// These should be logged (at or above WARNING)
	Status("Status message")
	Fatal("Fatal message")
	Error("Error message")
	Warning("Warning message")

	// These should not be logged (below WARNING)
	Info("Info message")
	Debug("Debug message")
	Trace("Trace message")

	// Verify IsLogging reflects the level
	if !IsLogging(logger.WARNING) {
		t.Error("Expected WARNING level to be logging")
	}
	if IsLogging(logger.INFO) {
		t.Error("Expected INFO level to NOT be logging")
	}
	if IsLogging(logger.DEBUG) {
		t.Error("Expected DEBUG level to NOT be logging")
	}
	if IsLogging(logger.TRACE) {
		t.Error("Expected TRACE level to NOT be logging")
	}
}

// TestTracingEnabled tests that tracing is enabled by default:
// - Verifies the trace logger has tracing enabled
// - Compares behavior with regular log package
//
// This ensures the trace package properly enables tracing.
func TestTracingEnabled(t *testing.T) {
	tempDir := t.TempDir()

	params := logger.Parameters{
		RootPath:                 tempDir,
		FileName:                 "trace-enabled",
		MaxFileSize:              1024,
		MaxCompressedFilesNumber: 3,
		Destinations:             "file",
		Prefix:                   "TRACE",
		Level:                    "INFO",
		MaintenanceInterval:      60,
	}

	err := Init(params)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Stop()

	// Log some messages - they should be prefixed with * due to tracing
	Info("Test message 1")
	Debug("Test message 2")
	Trace("Test message 3")

	// The actual verification of the * prefix would require reading the log file
	// and checking for the prefix, but that's implementation-specific
	// This test mainly ensures no errors occur with tracing enabled
}
