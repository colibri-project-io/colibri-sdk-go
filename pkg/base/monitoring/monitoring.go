package monitoring

import (
	"context"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
)

type monitoring interface {
	wrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request))
	startTransaction(name string) (interface{}, context.Context)
	endTransaction(transaction interface{})
	startTransactionSegment(transaction interface{}, name string, atributes map[string]interface{}) interface{}
	endTransactionSegment(segment interface{})
	getTransactionInContext(ctx context.Context) interface{}
	noticeError(transaction interface{}, err error)
}

var instance monitoring

func Initialize() {
	if config.IsProductionEnvironment() {
		instance = newProduction()
	} else {
		instance = newOthers()
	}
}

func WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return instance.wrapHandleFunc(pattern, handler)
}

func StartTransaction(name string) (interface{}, context.Context) {
	return instance.startTransaction(name)
}

func EndTransaction(transaction interface{}) {
	instance.endTransaction(transaction)
}

func StartTransactionSegment(transaction interface{}, name string, atributes map[string]interface{}) interface{} {
	return instance.startTransactionSegment(transaction, name, atributes)
}

func EndTransactionSegment(segment interface{}) {
	instance.endTransactionSegment(segment)
}

func GetTransactionInContext(ctx context.Context) interface{} {
	return instance.getTransactionInContext(ctx)
}

func NoticeError(transaction interface{}, err error) {
	instance.noticeError(transaction, err)
}
