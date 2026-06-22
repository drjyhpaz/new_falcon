package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorCyan    = "\033[36m"
)

var logFile *os.File

// Init initializes the logger
func Init(filename string) error {
	var err error
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	return err
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func logMessage(level, color, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	output := fmt.Sprintf("[%s] [%s%s%s] %s", timestamp, color, level, ColorReset, message)

	fmt.Println(output)

	if logFile != nil {
		logFile.WriteString(fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message))
	}
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	logMessage("INFO", ColorBlue, format, args...)
}

// Success logs a success message
func Success(format string, args ...interface{}) {
	logMessage("SUCCESS", ColorGreen, format, args...)
}

// Warning logs a warning message
func Warning(format string, args ...interface{}) {
	logMessage("WARNING", ColorYellow, format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	logMessage("ERROR", ColorRed, format, args...)
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	logMessage("DEBUG", ColorCyan, format, args...)
}
