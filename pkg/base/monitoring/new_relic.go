package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type newRelic struct {
	*newrelic.Application
}

func startNewRelicMonitoring() monitoring {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(fmt.Sprintf("%s-nr", config.APP_NAME)),
		newrelic.ConfigLicense(config.NEW_RELIC_LICENSE),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		logging.Fatal("An error occurred while loading the monitoring provider. Error: %s", err)
	}

	return &newRelic{app}
}

func (m *newRelic) startTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	transaction := m.Application.StartTransaction(name)
	nrctx := newrelic.NewContext(ctx, transaction)

	return transaction, nrctx
}

func (m *newRelic) endTransaction(transaction interface{}) {
	transaction.(*newrelic.Transaction).End()
}

func (m *newRelic) setWebRequest(_ context.Context, transaction interface{}, header http.Header, url *url.URL, method string) {
	transaction.(*newrelic.Transaction).SetWebRequest(newrelic.WebRequest{
		Header:    header,
		URL:       url,
		Method:    method,
		Transport: newrelic.TransportHTTP,
	})
}

func (m *newRelic) startWebRequest(ctx context.Context, header http.Header, path string, method string) (interface{}, context.Context) {
	txn, ctx := m.startTransaction(ctx, fmt.Sprintf("%s %s", method, path))
	m.setWebRequest(ctx, txn, header, &url.URL{Path: path}, method)
	return txn, ctx
}

func (m *newRelic) setWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter {
	return transaction.(*newrelic.Transaction).SetWebResponse(w)
}

func (m *newRelic) startTransactionSegment(_ context.Context, transaction interface{}, name string, attributes map[string]interface{}) interface{} {
	segment := transaction.(*newrelic.Transaction).StartSegment(name)
	segment.StartTime = transaction.(*newrelic.Transaction).StartSegmentNow()

	for key, value := range attributes {
		segment.AddAttribute(key, value)
	}

	return segment
}

func (m *newRelic) endTransactionSegment(segment interface{}) {
	segment.(*newrelic.Segment).End()
}

func (m *newRelic) getTransactionInContext(ctx context.Context) interface{} {
	return newrelic.FromContext(ctx)
}

func (m *newRelic) noticeError(transaction interface{}, err error) {
	transaction.(*newrelic.Transaction).NoticeError(err)
}
