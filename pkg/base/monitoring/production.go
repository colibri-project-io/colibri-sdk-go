package monitoring

import (
	"context"
	"net/http"
	"net/url"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type production struct {
	*newrelic.Application
}

func newProduction() *production {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(config.APP_NAME),
		newrelic.ConfigLicense(config.NEW_RELIC_LICENSE),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		logging.Fatal("An error occurred while loading the monitoring provider. Error: %s", err)
	}

	return &production{app}
}

func (m *production) wrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return newrelic.WrapHandleFunc(m.Application, pattern, handler)
}

func (m *production) startTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	transaction := m.Application.StartTransaction(name)
	nrctx := newrelic.NewContext(ctx, transaction)

	return transaction, nrctx
}

func (m *production) endTransaction(transaction interface{}) {
	transaction.(*newrelic.Transaction).End()
}

func (m *production) setWebRequest(transaction interface{}, header http.Header, url *url.URL, method string) {
	transaction.(*newrelic.Transaction).SetWebRequest(newrelic.WebRequest{
		Header:    header,
		URL:       url,
		Method:    method,
		Transport: newrelic.TransportHTTP,
	})
}

func (m *production) setWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter {
	return transaction.(*newrelic.Transaction).SetWebResponse(w)
}

func (m *production) startTransactionSegment(transaction interface{}, name string, atributes map[string]interface{}) interface{} {
	segment := transaction.(*newrelic.Transaction).StartSegment(name)
	segment.StartTime = transaction.(*newrelic.Transaction).StartSegmentNow()

	for key, value := range atributes {
		segment.AddAttribute(key, value)
	}

	return segment
}

func (m *production) endTransactionSegment(segment interface{}) {
	segment.(*newrelic.Segment).End()
}

func (m *production) getTransactionInContext(ctx context.Context) interface{} {
	return newrelic.FromContext(ctx)
}

func (m *production) noticeError(transaction interface{}, err error) {
	transaction.(*newrelic.Transaction).NoticeError(err)
}
