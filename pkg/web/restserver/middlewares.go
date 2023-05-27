package restserver

import (
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
)

const (
	userIDHeader   = "X-UserId"
	tenantIDHeader = "X-TenantId"
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

func authenticationContextFiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !strings.Contains(c.Request().URI().String(), string(AuthenticatedApi)) {
			return c.Next()
		}

		tenantID := extractUuidFromHeader(c, tenantIDHeader)
		userID := extractUuidFromHeader(c, userIDHeader)
		authCtx := security.NewAuthenticationContext(tenantID, userID)
		if authCtx.Valid() {
			newCtx := authCtx.SetInContext(c.UserContext())
			c.SetUserContext(newCtx)
			return c.Next()
		}

		c.Status(http.StatusUnauthorized)
		c.Request()
		return c.JSON(&Error{Error: "user not authenticated"})
	}
}

func extractUuidFromHeader(ctx *fiber.Ctx, key string) uuid.UUID {
	valueStr := string(ctx.Request().Header.Peek(key))
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

func newRelicFiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		txn, ctx := monitoring.StartTransaction(c.UserContext(), fmt.Sprintf("%s %s", c.Request().Header.Method(), c.Request().URI().Path()))
		defer monitoring.EndTransaction(txn)

		headers := make(http.Header)
		c.Context().Request.Header.VisitAll(func(key, value []byte) {
			headers.Set(string(key), string(value))
		})
		monitoring.SetWebRequest(ctx, txn, headers, &url.URL{Path: c.BaseURL()}, c.Method())

		c.SetUserContext(ctx)
		err := c.Next()

		if err != nil {
			monitoring.NoticeError(txn, err)
			return err
		}
		return nil
	}
}

func accessControlFiberMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "OPTIONS, GET, POST, PUT, PATCH, DELETE",
		AllowHeaders: "Origin, Content-Type",
	})
}
