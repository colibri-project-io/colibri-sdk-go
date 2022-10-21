package monitoring

import (
	"context"
	"net/http"

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

func (m *others) startTransaction(name string) (interface{}, context.Context) {
	logging.Info("Starting transaction monitoring with name %s", name)

	return nil, context.Background()
}

func (m *others) endTransaction(_ interface{}) {
	logging.Info("Ending transaction monitoring")
}

func (m *others) startTransactionSegment(_ interface{}, name string, _ map[string]interface{}) interface{} {
	logging.Info("Starting transaction segment monitoring with name %s", name)

	return nil
}

func (m *others) endTransactionSegment(_ interface{}) {
	logging.Info("Ending transaction segment monitoring")
}

func (m *others) getTransactionInContext(_ context.Context) interface{} {
	logging.Info("Getting transaction in context")

	return nil
}

func (m *others) noticeError(_ interface{}, err error) {
	logging.Info("Warning error %v", err)
}
