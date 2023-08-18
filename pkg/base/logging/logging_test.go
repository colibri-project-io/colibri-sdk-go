package logging

import (
	"bytes"
	"encoding/json"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func captureOutput(f func()) (out map[string]any) {
	var buf bytes.Buffer
	config.LOG_OUTPUT = &buf
	CreateLogger()
	f()
	_ = json.Unmarshal(buf.Bytes(), &out)
	return
}

func TestLogging(t *testing.T) {
	config.ENVIRONMENT = "test"
	config.APP_NAME = "sdk-test"
	config.APP_TYPE = "service"

	t.Run("Should log info", func(t *testing.T) {
		text := "Log info test"

		output := captureOutput(func() {
			Info(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "INFO", output["level"])
	})

	t.Run("Should log fatal", func(t *testing.T) {
		text := "Log fatal test"

		assert.PanicsWithValue(t, text, func() { Fatal(text) })
	})

	t.Run("Should log error", func(t *testing.T) {
		text := "Log error test"

		output := captureOutput(func() {
			Error(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "ERROR", output["level"])
	})

	t.Run("Should log warn", func(t *testing.T) {
		text := "Log warn test"

		output := captureOutput(func() {
			Warn(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "WARN", output["level"])
	})
}

func TestDebug(t *testing.T) {
	config.APP_NAME = "sdk-test"
	config.APP_TYPE = "service"
	config.LOG_LEVEL = "debug"

	_ = config.Load()

	t.Run("Should log debug", func(t *testing.T) {
		text := "Log debug test"

		output := captureOutput(func() {
			Debug(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "DEBUG", output["level"])
	})
}
