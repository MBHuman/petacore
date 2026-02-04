package logger

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Init initializes the global logger
func Init(development bool) error {
	var err error
	var config zap.Config

	if development {
		config = zap.NewDevelopmentConfig()
		config.DisableStacktrace = true
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err = config.Build()
	if err != nil {
		return err
	}

	sugar = logger.Sugar()
	return nil
}

func SetLevel(level zap.AtomicLevel) {
	if logger != nil {
		logger = logger.WithOptions(zap.IncreaseLevel(level))
		sugar = logger.Sugar()
	}
}

// GetLogger returns the underlying zap logger
func GetLogger() *zap.Logger {
	if logger == nil {
		// Fallback to nop logger
		logger = zap.NewNop()
	}
	return logger
}

// GetSugar returns the sugared logger
func GetSugar() *zap.SugaredLogger {
	if sugar == nil {
		sugar = GetLogger().Sugar()
	}
	return sugar
}

// Sync flushes any buffered log entries
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Debugf logs a debug message with fmt-style formatting
func Debugf(template string, args ...interface{}) {
	GetSugar().Debugf(template, args...)
}

// Infof logs an info message with fmt-style formatting
func Infof(template string, args ...interface{}) {
	GetSugar().Infof(template, args...)
}

// Warnf logs a warning message with fmt-style formatting
func Warnf(template string, args ...interface{}) {
	GetSugar().Warnf(template, args...)
}

// Errorf logs an error message with fmt-style formatting
func Errorf(template string, args ...interface{}) {
	GetSugar().Errorf(template, args...)
}

// Fatalf logs a fatal message with fmt-style formatting and exits
func Fatalf(template string, args ...interface{}) {
	GetSugar().Fatalf(template, args...)
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}
