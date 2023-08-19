package logging

import (
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"log/slog"
	"strings"
)

var log *slog.Logger

func init() {
	CreateLogger()
}

func CreateLogger() {
	opts := &slog.HandlerOptions{Level: parseLevel(config.LOG_LEVEL)}
	log = slog.New(slog.NewJSONHandler(config.LOG_OUTPUT, opts)).
		With("colibri-sdk-go", map[string]string{"application": config.APP_NAME, "application_type": config.APP_TYPE, "version": config.VERSION})
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
	log.Info(msg)
}

// Fatal prints in console a fatal message and exits program
func Fatal(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Error(msg)
	panic(msg)
}

// Error prints in console a error message
func Error(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Error(msg)
}

// Warn prints in console a warn message
func Warn(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Warn(msg)
}

// Debug prints in console a debug message if config is enabled
func Debug(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Debug(msg)
}
