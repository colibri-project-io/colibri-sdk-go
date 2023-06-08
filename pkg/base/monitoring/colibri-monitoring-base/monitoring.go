package colibri_monitoring_base

import (
	"context"
	"net/http"
)

// Monitoring is a contract to implements all necessary functions
type Monitoring interface {
	StartTransaction(ctx context.Context, name string) (interface{}, context.Context)
	EndTransaction(transaction interface{})
	StartWebRequest(ctx context.Context, header http.Header, path string, method string) (interface{}, context.Context)
	StartTransactionSegment(ctx context.Context, name string, attributes map[string]string) interface{}
	EndTransactionSegment(segment interface{})
	GetTransactionInContext(ctx context.Context) interface{}
	NoticeError(transaction interface{}, err error)
	GetSQLDBDriverName() string
}
