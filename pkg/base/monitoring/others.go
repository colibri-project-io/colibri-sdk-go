package monitoring

import (
	"context"
	"net/http"
	"net/url"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type others struct {
}

func newOthers() *others {
	return &others{}
}

func (m *others) wrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return pattern, handler
}

func (m *others) startTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	logging.Debug("Starting transaction monitoring with name %s", name)
	return nil, ctx
}

func (m *others) endTransaction(_ interface{}) {
	logging.Debug("Ending transaction monitoring")
}

func (m *others) setWebRequest(transaction interface{}, header http.Header, url *url.URL, method string) {
	logging.Debug("Setting web request in transaction with path %s", url.Path)
}

func (m *others) setWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter {
	logging.Debug("Setting web response in transaction")
	return w
}

func (m *others) startTransactionSegment(_ interface{}, name string, _ map[string]interface{}) interface{} {
	logging.Debug("Starting transaction segment monitoring with name %s", name)
	return nil
}

func (m *others) endTransactionSegment(_ interface{}) {
	logging.Debug("Ending transaction segment monitoring")
}

func (m *others) getTransactionInContext(_ context.Context) interface{} {
	logging.Debug("Getting transaction in context")
	return nil
}

func (m *others) noticeError(_ interface{}, err error) {
	logging.Debug("Warning error %v", err)
}
