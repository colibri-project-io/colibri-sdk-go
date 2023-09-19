package logging_test

import (
	"bytes"
	"encoding/json"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func captureOutput(f func()) (out map[string]any) {
	var buf bytes.Buffer
	config.LOG_OUTPUT = &buf
	logging.InitializeLogger()
	f()
	_ = json.Unmarshal(buf.Bytes(), &out)
	return
}

func TestLogging(t *testing.T) {
	config.ENVIRONMENT = "test"
	config.APP_NAME = "sdk-test"
	config.APP_TYPE = "service"

	t.Run("Should logger info", func(t *testing.T) {
		text := "Log info test"

		output := captureOutput(func() {
			logging.Info(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "INFO", output["level"])
	})

	t.Run("Should logger fatal", func(t *testing.T) {
		text := "Log fatal test"

		assert.PanicsWithValue(t, text, func() { logging.Fatal(text) })
	})

	t.Run("Should logger error", func(t *testing.T) {
		text := "Log error test"

		output := captureOutput(func() {
			logging.Error(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "ERROR", output["level"])
	})

	t.Run("Should logger warn", func(t *testing.T) {
		text := "Log warn test"

		output := captureOutput(func() {
			logging.Warn(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "WARN", output["level"])
	})
}

func TestDebug(t *testing.T) {
	config.APP_NAME = "sdk-test"
	config.APP_TYPE = "service"
	config.LOG_LEVEL = "debug"

	logging.InitializeLogger()

	t.Run("Should logger debug", func(t *testing.T) {
		text := "Log debug test"

		output := captureOutput(func() {
			logging.Debug(text)
		})

		assert.Equal(t, text, output["msg"])
		assert.Equal(t, "DEBUG", output["level"])
	})
}
