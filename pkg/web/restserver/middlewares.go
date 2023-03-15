package restserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/google/uuid"
)

const (
	userIDHeader      = "X-UserId"
	tenantIDHeader    = "X-TenantId"
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON   = "application/json"
)

type MiddlewareError struct {
	Err        error `json:"error"`
	StatusCode int   `json:"statusCode"`
}

func (e MiddlewareError) Error() string {
	return e.Err.Error()
}

func NewMiddlewareError(statusCode int, err error) *MiddlewareError {
	return &MiddlewareError{StatusCode: statusCode, Err: err}
}

type CustomMiddleware interface {
	Apply(ctx WebContext) *MiddlewareError
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authenticationContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, string(AuthenticatedApi)) {
			next.ServeHTTP(w, r)
			return
		}

		tenantID := extractUuidFromHeader(r, tenantIDHeader)
		userID := extractUuidFromHeader(r, userIDHeader)
		authCtx := security.NewAuthenticationContext(tenantID, userID)
		if authCtx.Valid() {
			next.ServeHTTP(w, r.WithContext(authCtx.SetInContext(r.Context())))
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(errors.New("user not authenticated")); err != nil {
			logging.Error(err.Error())
		}
	})
}

func extractUuidFromHeader(r *http.Request, key string) uuid.UUID {
	valueStr := r.Header.Get(key)
	if valueStr == "" {
		return uuid.Nil
	}

	value, err := uuid.Parse(valueStr)
	if err != nil {
		logging.Error("could not parse %s from header %s: %v", key, valueStr, err)
		return uuid.Nil
	}

	return value
}

func newRelicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn, ctx := monitoring.StartTransaction(r.Context(), r.Method+" "+r.URL.Path)
		defer monitoring.EndTransaction(txn)

		monitoring.SetWebRequest(txn, r.Header, r.URL, r.Method)
		w = monitoring.SetWebResponse(txn, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
