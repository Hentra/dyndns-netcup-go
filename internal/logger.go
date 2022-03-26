package internal

import "log"

// Logger represents an logger instance
type Logger struct {
	verbose bool
}

// NewLogger creates a Logger instance with given verboseness
func NewLogger(verbose bool) *Logger {
	return &Logger{verbose}
}

// Info logs a info message which will only show when verbose is set
func (l *Logger) Info(msg string, v ...interface{}) {
	if l.verbose {
		log.Printf(msg, v...)
	}
}

// Warning will log a warning message
func (l *Logger) Warning(msg string, v ...interface{}) {
	log.Printf("[Warning]: "+msg, v...)
}

// Error will log an error message and exit
func (l *Logger) Error(v ...interface{}) {
	log.Fatal(v...)
}
