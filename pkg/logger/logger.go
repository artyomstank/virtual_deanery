// pkg/logger/logger.go
package logger

import (
	"context"

	"go.uber.org/zap"
)

// Logger provides logging interface.
type Logger interface {
	// Debug logs at debug level.
	Debug(msg string, fields ...interface{})

	// Info logs at info level.
	Info(msg string, fields ...interface{})

	// Warn logs at warning level.
	Warn(msg string, fields ...interface{})

	// Error logs at error level.
	Error(msg string, fields ...interface{})

	// Fatal logs at fatal level and exits.
	Fatal(msg string, fields ...interface{})

	// WithContext returns context-aware logger.
	WithContext(ctx context.Context) Logger
}

// ZapLogger wraps zap logger.
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates new zap logger.
func NewZapLogger(level string) (*ZapLogger, error) {
	// TODO: Create zap config based on level (debug, info, warn, error)

	// TODO: Build logger from config

	// TODO: Return ZapLogger or error
	return nil, nil
}

// Debug logs at debug level.
func (l *ZapLogger) Debug(msg string, fields ...interface{}) {
	// TODO: Convert fields to zap fields (if any)

	// TODO: Log using l.logger.Debug()
}

// Info logs at info level.
func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	// TODO: Log using l.logger.Info()
}

// Warn logs at warning level.
func (l *ZapLogger) Warn(msg string, fields ...interface{}) {
	// TODO: Log using l.logger.Warn()
}

// Error logs at error level.
func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	// TODO: Log using l.logger.Error()
}

// Fatal logs at fatal level and exits.
func (l *ZapLogger) Fatal(msg string, fields ...interface{}) {
	// TODO: Log using l.logger.Fatal()
}

// WithContext returns context-aware logger.
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	// TODO: Return logger that includes context info in logs
	return l
}

// Sync flushes any buffered log entries.
func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
