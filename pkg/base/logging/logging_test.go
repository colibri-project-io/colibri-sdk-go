package logging

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func formatExpected(level Level, text string) string {
	return fmt.Sprintf("%s %s %s\n", time.Now().Format("2006/01/02 15:04:05"), level, text)
}

func TestLogging(t *testing.T) {
	t.Run("Should log info", func(t *testing.T) {
		text := "Log info test"

		output := captureOutput(func() {
			Info(text)
		})

		assert.Equal(t, formatExpected(INFO, text), output)
	})

	t.Run("Should log error", func(t *testing.T) {
		text := "Log error test"

		output := captureOutput(func() {
			Error(text)
		})

		assert.Equal(t, formatExpected(ERROR, text), output)
	})

	t.Run("Should log warn", func(t *testing.T) {
		text := "Log warn test"

		output := captureOutput(func() {
			Warn(text)
		})

		assert.Equal(t, formatExpected(WARN, text), output)
	})
}
