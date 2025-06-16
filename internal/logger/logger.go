package logger

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Logger defines minimal logging interface for our application.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	// New: Add Warn method to the interface
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
}

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// defaultLogger is a simple thread-safe console logger.
type defaultLogger struct {
	mu sync.Mutex
}

// New creates a new defaultLogger instance.
func New() Logger {
	return &defaultLogger{}
}

// Info logs an informational message.
func (l *defaultLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.print("INFO", msg, nil, fields...)
}

// New: Implement the Warn method for defaultLogger
// Warn logs a warning message.
func (l *defaultLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.print("WARN", msg, nil, fields...) // Warnings usually don't have an 'err' object
}

// Error logs an error message.
func (l *defaultLogger) Error(ctx context.Context, msg string, err error, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.print("ERROR", msg, err, fields...)
}

// print formats and outputs the log entry.
func (l *defaultLogger) print(level, msg string, err error, fields ...Field) {
	timestamp := time.Now().Format(time.RFC3339)
	entry := fmt.Sprintf("%s [%s] %s", timestamp, level, msg)
	if err != nil {
		entry += fmt.Sprintf(" error=%v", err)
	}
	for _, f := range fields {
		entry += fmt.Sprintf(" %s=%v", f.Key, f.Value)
	}
	fmt.Println(entry)
}
