package monitoring

import (
	"context"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"net/http"
	"net/url"
)

// monitoring is a contract to implements all necessary functions
type monitoring interface {
	startTransaction(ctx context.Context, name string) (interface{}, context.Context)
	endTransaction(transaction interface{})
	setWebRequest(ctx context.Context, transaction interface{}, header http.Header, url *url.URL, method string)
	startWebRequest(ctx context.Context, header http.Header, path string, method string) (interface{}, context.Context)
	setWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter
	startTransactionSegment(ctx context.Context, transaction interface{}, name string, attributes map[string]interface{}) interface{}
	endTransactionSegment(segment interface{})
	getTransactionInContext(ctx context.Context) interface{}
	noticeError(transaction interface{}, err error)
}

var instance monitoring

// Initialize loads the monitoring settings according to the configured environment.
func Initialize() {
	if useNRMonitoring() {
		instance = startNewRelicMonitoring()
	} else if useOTELMonitoring() {
		instance = startOpenTelemetryMonitoring()
	} else {
		instance = newOthers()
	}
}

func useOTELMonitoring() bool {
	return config.OTEL_EXPORTER_OTLP_ENDPOINT != ""
}

func useNRMonitoring() bool {
	return config.NEW_RELIC_LICENSE != ""
}

// StartTransaction start a transaction in context with name
func StartTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	return instance.startTransaction(ctx, name)
}

// EndTransaction ends the transaction
func EndTransaction(transaction interface{}) {
	instance.endTransaction(transaction)
}

// StartWebRequest sets a web request config inside transaction
func StartWebRequest(ctx context.Context, header http.Header, path string, method string) (interface{}, context.Context) {
	return instance.startWebRequest(ctx, header, path, method)
}

// SetWebResponse sets a web response config inside transaction TODO Is this still used?
func SetWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter {
	return instance.setWebResponse(transaction, w)
}

// StartTransactionSegment start a transaction segment inside opened transaction with name and atributes
func StartTransactionSegment(ctx context.Context, transaction interface{}, name string, attributes map[string]interface{}) interface{} {
	return instance.startTransactionSegment(ctx, transaction, name, attributes)
}

// EndTransactionSegment ends the transaction segment
func EndTransactionSegment(segment interface{}) {
	instance.endTransactionSegment(segment)
}

// GetTransactionInContext returns transaction inside a context
func GetTransactionInContext(ctx context.Context) interface{} {
	return instance.getTransactionInContext(ctx)
}

// NoticeError notices an error in monitoring provider
func NoticeError(transaction interface{}, err error) {
	instance.noticeError(transaction, err)
}
