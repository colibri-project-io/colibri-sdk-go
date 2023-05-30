package colibri_monitoring_base

import (
	"context"
	"net/http"
	"net/url"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type others struct {
}

func NewOthers() Monitoring {
	return &others{}
}

func (m *others) StartTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	logging.Debug("Starting transaction Monitoring with name %s", name)
	return nil, ctx
}

func (m *others) EndTransaction(_ interface{}) {
	logging.Debug("Ending transaction Monitoring")
}

func (m *others) SetWebRequest(_ context.Context, transaction interface{}, header http.Header, url *url.URL, method string) {
	logging.Debug("Setting web request in transaction with path %s", url.Path)
}

func (m *others) StartWebRequest(ctx context.Context, header http.Header, path string, method string) (interface{}, context.Context) {
	logging.Debug("Start web request in transaction with path %s", path)
	return nil, ctx
}

func (m *others) SetWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter {
	logging.Debug("Setting web response in transaction")
	return w
}

func (m *others) StartTransactionSegment(_ context.Context, name string, _ map[string]string) interface{} {
	logging.Debug("Starting transaction segment Monitoring with name %s", name)
	return nil
}

func (m *others) EndTransactionSegment(_ interface{}) {
	logging.Debug("Ending transaction segment Monitoring")
}

func (m *others) GetTransactionInContext(_ context.Context) interface{} {
	logging.Debug("Getting transaction in context")
	return nil
}

func (m *others) NoticeError(_ interface{}, err error) {
	logging.Debug("Warning error %v", err)
}

func (m *others) GetSQLDBDriverName() string {
	return "postgres"
}
