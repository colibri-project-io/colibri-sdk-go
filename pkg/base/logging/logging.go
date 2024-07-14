package logging

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
)

var logger *slog.Logger

func init() {
	InitializeLogger()
}

func InitializeLogger() {
	logger = slog.New(createLogHandler())
}

func createLogHandler() slog.Handler {
	opts := &slog.HandlerOptions{Level: parseLevel(config.LOG_LEVEL)}
	if config.IsDevelopmentEnvironment() {
		return slog.NewTextHandler(config.LOG_OUTPUT, opts)
	}
	return slog.NewJSONHandler(config.LOG_OUTPUT, opts)
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "error":
		return slog.LevelError
	case "warn", "warning":
		return slog.LevelWarn
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// Info prints in console an info message
func Info(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	logger.Info(msg)
}

// Fatal prints in console a fatal message and exits program
func Fatal(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	logger.Error(msg)
	panic(msg)
}

// Error prints in console a error message
func Error(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	logger.Error(msg)
}

// Warn prints in console a warn message
func Warn(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	logger.Warn(msg)
}

// Debug prints in console a debug message if config is enabled
func Debug(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	logger.Debug(msg)
}
