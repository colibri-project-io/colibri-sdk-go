package logging

import (
	"fmt"
	"log"
)

type Level string

const (
	INFO        Level = "INFO"
	FATAL       Level = "FATAL"
	ERROR       Level = "ERROR"
	WARN        Level = "WARN"
	DEFAULT_MSG       = "%s %s"
)

func Info(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Printf(DEFAULT_MSG, INFO, msg)
}

func Fatal(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Fatalf(DEFAULT_MSG, FATAL, msg)
}

func Error(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Printf(DEFAULT_MSG, ERROR, msg)
}

func Warn(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	log.Printf(DEFAULT_MSG, WARN, msg)
}
