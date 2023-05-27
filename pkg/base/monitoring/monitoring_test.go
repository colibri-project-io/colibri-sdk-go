package monitoring

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/stretchr/testify/assert"
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
		txnName := "txn-test"
		segName := "txn-segment-test"
		var w http.ResponseWriter

		transaction, ctx := StartTransaction(context.Background(), txnName)
		SetWebRequest(ctx, transaction, http.Header{}, &url.URL{}, http.MethodGet)
		SetWebResponse(transaction, w)
		segment := StartTransactionSegment(ctx, transaction, segName, map[string]interface{}{
			"TestKey": "TestValue",
		})

		EndTransactionSegment(segment)
		NoticeError(transaction, errors.New("Test notice error"))
		EndTransaction(transaction)

		assert.NotNil(t, transaction)
		assert.NotNil(t, segment)
		assert.NotEmpty(t, ctx)
	})
}

func TestNonProductionMonitoring(t *testing.T) {
	captureOutput := func(fn func()) string {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		fn()
		log.SetOutput(os.Stderr)
		return buf.String()
	}

	formatExpected := func(text string) string {
		return fmt.Sprintf("%s DEBUG %s\n", time.Now().Format("2006/01/02 15:04:05"), text)
	}

	config.APP_NAME = "colibri-project-test"
	config.ENVIRONMENT = config.ENVIRONMENT_TEST
	config.DEBUG = true
	assert.False(t, config.IsProductionEnvironment())

	Initialize()
	assert.NotNil(t, instance)

	t.Run("Should start transaction", func(t *testing.T) {
		name := "txn-test"
		text := fmt.Sprintf("Starting transaction monitoring with name %s", name)

		output := captureOutput(func() {
			transaction, ctx := StartTransaction(context.Background(), name)
			assert.Nil(t, transaction)
			assert.Empty(t, ctx)
		})

		assert.Equal(t, formatExpected(text), output)
	})

	t.Run("Should end transaction", func(t *testing.T) {
		text := "Ending transaction monitoring"

		output := captureOutput(func() {
			EndTransaction(text)
		})

		assert.Equal(t, formatExpected(text), output)
	})

	t.Run("Should start transaction segment", func(t *testing.T) {
		name := "txn-segment-test"
		text := fmt.Sprintf("Starting transaction segment monitoring with name %s", name)

		output := captureOutput(func() {
			segment := StartTransactionSegment(context.Background(), name, name, nil)
			assert.Nil(t, segment)
		})

		assert.Equal(t, formatExpected(text), output)
	})

	t.Run("Should end transaction segment", func(t *testing.T) {
		text := "Ending transaction segment monitoring"

		output := captureOutput(func() {
			EndTransactionSegment(text)
		})

		assert.Equal(t, formatExpected(text), output)
	})

	t.Run("Should get transaction in context", func(t *testing.T) {
		text := "Getting transaction in context"

		output := captureOutput(func() {
			GetTransactionInContext(context.Background())
		})

		assert.Equal(t, formatExpected(text), output)
	})

	t.Run("Should notice error", func(t *testing.T) {
		err := errors.New("error test")
		text := fmt.Sprintf("Warning error %v", err)

		output := captureOutput(func() {
			NoticeError(text, err)
		})

		assert.Equal(t, formatExpected(text), output)
	})
}
