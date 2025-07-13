package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logLevel = INFO
	logger   = log.New(os.Stdout, "", 0)
)

func init() {
	// Set default log flags to include timestamp
	logger.SetFlags(0)
}

// SetLevel sets the minimum log level
func SetLevel(level LogLevel) {
	logLevel = level
}

// SetLevelFromString sets log level from string
func SetLevelFromString(level string) {
	switch strings.ToLower(level) {
	case "debug":
		SetLevel(DEBUG)
	case "info":
		SetLevel(INFO)
	case "warn", "warning":
		SetLevel(WARN)
	case "error":
		SetLevel(ERROR)
	case "fatal":
		SetLevel(FATAL)
	default:
		SetLevel(INFO)
	}
}

// getCallerInfo returns the file and line number of the caller
func getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0
	}
	
	// Get just the filename, not the full path
	filename := filepath.Base(file)
	return filename, line
}

// formatMessage formats the log message with timestamp, level, file, line, and message
func formatMessage(level string, message string, skip int) string {
	file, line := getCallerInfo(skip)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("[%s] %s %s:%d - %s", timestamp, level, file, line, message)
}

// logf is the internal logging function
func logf(level LogLevel, levelStr string, format string, args ...interface{}) {
	if level < logLevel {
		return
	}
	
	message := fmt.Sprintf(format, args...)
	formatted := formatMessage(levelStr, message, 3) // skip 3 frames: logf -> Debug/Info/etc -> caller
	logger.Println(formatted)
}

// Debug logs debug messages
func Debug(format string, args ...interface{}) {
	logf(DEBUG, "DEBUG", format, args...)
}

// Info logs info messages
func Info(format string, args ...interface{}) {
	logf(INFO, "INFO", format, args...)
}

// Warn logs warning messages
func Warn(format string, args ...interface{}) {
	logf(WARN, "WARN", format, args...)
}

// Error logs error messages
func Error(format string, args ...interface{}) {
	logf(ERROR, "ERROR", format, args...)
}

// Fatal logs fatal messages and exits
func Fatal(format string, args ...interface{}) {
	logf(FATAL, "FATAL", format, args...)
	os.Exit(1)
}

// Convenience functions for common use cases

// Printf logs an info message (alias for Info)
func Printf(format string, args ...interface{}) {
	Info(format, args...)
}

// Println logs an info message with automatic formatting
func Println(args ...interface{}) {
	Info(fmt.Sprint(args...))
}

// ErrorIf logs an error if err is not nil
func ErrorIf(err error, format string, args ...interface{}) {
	if err != nil {
		message := fmt.Sprintf(format, args...)
		Error("%s: %v", message, err)
	}
}

// DebugIf logs a debug message if condition is true
func DebugIf(condition bool, format string, args ...interface{}) {
	if condition {
		Debug(format, args...)
	}
}