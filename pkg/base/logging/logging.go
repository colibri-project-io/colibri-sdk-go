package logging

import (
	"fmt"
	"log"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
)

type Level string

const (
	INFO        Level = "INFO"
	FATAL       Level = "FATAL"
	ERROR       Level = "ERROR"
	WARN        Level = "WARN"
	DEBUG       Level = "DEBUG"
	DEFAULT_MSG       = "%s %s"
)

// Info prints in console a info message
func Info(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Printf(DEFAULT_MSG, INFO, msg)
}

// Fatal prints in console a fatal message and exits program
func Fatal(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Printf(DEFAULT_MSG, FATAL, msg)
	panic(msg)
}

// Error prints in console a error message
func Error(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Printf(DEFAULT_MSG, ERROR, msg)
}

// Warn prints in console a warn message
func Warn(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Printf(DEFAULT_MSG, WARN, msg)
}

// Debug prints in console a debug message if config is enabled
func Debug(message string, args ...interface{}) {
	if config.DEBUG {
		msg := fmt.Sprintf(message, args...)
		log.Printf(DEFAULT_MSG, DEBUG, msg)
	}
}
