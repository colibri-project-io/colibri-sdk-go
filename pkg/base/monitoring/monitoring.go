package monitoring

import (
	"context"
	"net/http"
	"net/url"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
)

// monitoring is a contract to implements all necessary functions
type monitoring interface {
	wrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request))
	startTransaction(ctx context.Context, name string) (interface{}, context.Context)
	endTransaction(transaction interface{})
	setWebRequest(transaction interface{}, header http.Header, url *url.URL, method string)
	setWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter
	startTransactionSegment(transaction interface{}, name string, atributes map[string]interface{}) interface{}
	endTransactionSegment(segment interface{})
	getTransactionInContext(ctx context.Context) interface{}
	noticeError(transaction interface{}, err error)
}

var instance monitoring

// Initialize loads the monitoring settings according to the configured environment.
func Initialize() {
	if config.IsProductionEnvironment() {
		instance = newProduction()
	} else {
		instance = newOthers()
	}
}

// WrapHandleFunc wrap the http rest handle functions.
func WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return instance.wrapHandleFunc(pattern, handler)
}

// StartTransaction start a transaction in context with name
func StartTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	return instance.startTransaction(ctx, name)
}

// EndTransaction ends the transaction
func EndTransaction(transaction interface{}) {
	instance.endTransaction(transaction)
}

// SetWebRequest sets a web request config inside transaction
func SetWebRequest(transaction interface{}, header http.Header, url *url.URL, method string) {
	instance.setWebRequest(transaction, header, url, method)
}

// SetWebResponse sets a web response config inside transaction
func SetWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter {
	return instance.setWebResponse(transaction, w)
}

// StartTransactionSegment start a transaction segment inside opened transaction with name and atributes
func StartTransactionSegment(transaction interface{}, name string, atributes map[string]interface{}) interface{} {
	return instance.startTransactionSegment(transaction, name, atributes)
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
