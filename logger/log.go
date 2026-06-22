package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

var levelNames = map[LogLevel]string{
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARN",
	Error: "ERROR",
}

type Logger struct {
	file      *os.File
	level     LogLevel
	mu        sync.Mutex
	callbacks []LogCallback
}

type LogEntry struct {
	Time    time.Time
	Level   LogLevel
	Message string
}

type LogCallback func(*LogEntry)

// NewLogger creates a new logger instance
func NewLogger(filename string, level string) (*Logger, error) {
	logLevel := Info
	switch level {
	case "debug":
		logLevel = Debug
	case "warn":
		logLevel = Warn
	case "error":
		logLevel = Error
	}

	// Create logs directory if not exists
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		os.MkdirAll(dir, 0755)
	}

	// Open or create log file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file:  file,
		level: logLevel,
	}, nil
}

// AddCallback adds a callback function for log entries
func (l *Logger) AddCallback(cb LogCallback) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.callbacks = append(l.callbacks, cb)
}

// log writes log entry
func (l *Logger) log(level LogLevel, msg string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry := &LogEntry{
		Time:    now,
		Level:   level,
		Message: msg,
	}

	// Format log line
	timestamp := now.Format("2006-01-02 15:04:05")
	levels := levelNames[level]
	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, levels, msg)

	// Write to file
	if l.file != nil {
		l.file.WriteString(logLine)
	}

	// Call callbacks
	for _, cb := range l.callbacks {
		cb(entry)
	}
}

// Debug logs debug message
func (l *Logger) Debug(msg string) {
	l.log(Debug, msg)
}

// Debugf logs debug message with formatting
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(Debug, fmt.Sprintf(format, args...))
}

// Info logs info message
func (l *Logger) Info(msg string) {
	l.log(Info, msg)
}

// Infof logs info message with formatting
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(Info, fmt.Sprintf(format, args...))
}

// Warn logs warning message
func (l *Logger) Warn(msg string) {
	l.log(Warn, msg)
}

// Warnf logs warning message with formatting
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(Warn, fmt.Sprintf(format, args...))
}

// Error logs error message
func (l *Logger) Error(msg string) {
	l.log(Error, msg)
}

// Errorf logs error message with formatting
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(Error, fmt.Sprintf(format, args...))
}

// Success logs success message
func (l *Logger) Success(msg string) {
	l.log(Info, "[SUCCESS] "+msg)
}

// Successf logs success message with formatting
func (l *Logger) Successf(format string, args ...interface{}) {
	l.Success(fmt.Sprintf(format, args...))
}

// Close closes the logger
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
