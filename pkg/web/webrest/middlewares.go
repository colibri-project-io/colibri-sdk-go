package webrest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/google/uuid"
)

const (
	userIDHeader   = "X-UserId"
	tenantIDHeader = "X-TenantId"
)

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
		if tenantID != uuid.Nil && userID != uuid.Nil {
			next.ServeHTTP(w, r.WithContext(security.NewAuthenticationContext(tenantID, userID).SetInContext(r.Context())))
			return
		}

		ErrorResponse(r, w, http.StatusUnauthorized, errors.New("user not authenticated"))
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
