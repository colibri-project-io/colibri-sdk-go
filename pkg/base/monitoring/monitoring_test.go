package monitoring

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestProductionMonitoring(t *testing.T) {
	config.NEW_RELIC_LICENSE = "abcdefghijklmopqrstuvwxyz1234567890aNRAL"
	config.APP_NAME = "test"
	config.ENVIRONMENT = config.ENVIRONMENT_PRODUCTION
	assert.True(t, config.IsProductionEnvironment())

	Initialize()
	assert.NotNil(t, instance)

	t.Run("Should get transaction in context", func(t *testing.T) {
		txnName := "txn-test"

		_, ctx := StartTransaction(context.Background(), txnName)
		transaction := GetTransactionInContext(ctx)
		EndTransaction(transaction)

		assert.NotNil(t, transaction)
	})

	t.Run("Should get nil when transaction not in context", func(t *testing.T) {
		transaction := GetTransactionInContext(context.Background())

		assert.Nil(t, transaction)
	})

	t.Run("Should start/end transaction, start/end segment and notice error", func(t *testing.T) {
		segName := "txn-segment-test"

		transaction, ctx := StartWebRequest(context.Background(), http.Header{}, "/", http.MethodGet)
		segment := StartTransactionSegment(ctx, segName, map[string]string{
			"TestKey": "TestValue",
		})

		EndTransactionSegment(segment)
		NoticeError(transaction, errors.New("test notice error"))
		EndTransaction(transaction)

		assert.NotNil(t, transaction)
		assert.NotNil(t, segment)
		assert.NotEmpty(t, ctx)
	})
}

func TestNonProductionMonitoring(t *testing.T) {
	captureOutput := func(fn func()) (out map[string]any) {
		var buf bytes.Buffer
		config.LOG_OUTPUT = &buf
		logging.CreateLogger()
		fn()
		_ = json.Unmarshal(buf.Bytes(), &out)
		return
	}

	config.APP_NAME = "colibri-project-test"
	config.ENVIRONMENT = config.ENVIRONMENT_TEST
	config.NEW_RELIC_LICENSE = ""
	config.LOG_LEVEL = "debug"
	assert.False(t, config.IsProductionEnvironment())

	Initialize()
	assert.NotNil(t, instance)

	t.Run("Should start transaction", func(t *testing.T) {
		name := "txn-test"
		text := fmt.Sprintf("Starting transaction Monitoring with name %s", name)

		output := captureOutput(func() {
			transaction, ctx := StartTransaction(context.Background(), name)
			assert.Nil(t, transaction)
			assert.Empty(t, ctx)
		})

		assert.Equal(t, text, output["msg"])
	})

	t.Run("Should end transaction", func(t *testing.T) {
		text := "Ending transaction Monitoring"

		output := captureOutput(func() {
			EndTransaction(text)
		})

		assert.Equal(t, text, output["msg"])
	})

	t.Run("Should start transaction segment", func(t *testing.T) {
		name := "txn-segment-test"
		text := fmt.Sprintf("Starting transaction segment Monitoring with name %s", name)

		output := captureOutput(func() {
			segment := StartTransactionSegment(context.Background(), name, nil)
			assert.Nil(t, segment)
		})

		assert.Equal(t, text, output["msg"])
	})

	t.Run("Should end transaction segment", func(t *testing.T) {
		text := "Ending transaction segment Monitoring"

		output := captureOutput(func() {
			EndTransactionSegment(text)
		})

		assert.Equal(t, text, output["msg"])
	})

	t.Run("Should get transaction in context", func(t *testing.T) {
		text := "Getting transaction in context"

		output := captureOutput(func() {
			GetTransactionInContext(context.Background())
		})

		assert.Equal(t, text, output["msg"])
	})

	t.Run("Should notice error", func(t *testing.T) {
		err := errors.New("error test")
		text := fmt.Sprintf("Warning error %v", err)

		output := captureOutput(func() {
			NoticeError(text, err)
		})

		assert.Equal(t, text, output["msg"])
	})
}
