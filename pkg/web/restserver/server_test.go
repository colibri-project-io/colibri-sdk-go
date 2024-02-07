package restserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restclient"
	"github.com/stretchr/testify/assert"
)

const (
	UserId   = "6e859dae-c879-11eb-b8bc-0242ac130004"
	TenantId = "6e859dae-c879-11eb-b8bc-0242ac130005"
)

type customMiddlewareTest struct {
}

func (t *customMiddlewareTest) Apply(ctx WebContext) *MiddlewareError {
	if strings.HasPrefix(ctx.Path(), "/public/middleware/conflict") {
		return NewMiddlewareError(http.StatusConflict, errors.New("has conflict"))
	}
	return nil
}

type customAuthenticationContextMiddleware struct {
}

func (m *customAuthenticationContextMiddleware) Apply(ctx WebContext) (*security.AuthenticationContext, error) {
	if !strings.Contains(ctx.Path(), string(AuthenticatedApi)) {
		return nil, nil
	}

	authHeaders := ctx.RequestHeaders()["Authorization"]
	if len(authHeaders) < 1 || authHeaders[0] == "" {
		return nil, UserUnauthenticatedError
	}
	return security.NewAuthenticationContext("ab123", "123abc"), nil

}

func beforeEnterApply(ctx WebContext) *MiddlewareError {
	if strings.HasPrefix(ctx.Path(), "/public/hook/conflict") {
		return NewMiddlewareError(http.StatusConflict, errors.New("has hook conflict"))
	}
	return nil
}

func TestStartRestServer(t *testing.T) {
	test.InitializeBaseTest()

	type Resp struct {
		Msg string `json:"msg" validate:"required"`
	}

	type Query struct {
		Msg  string `form:"msg"`
		Size uint8  `form:"size"`
	}

	listener := func() (l net.Listener) {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Error("could not get available port")
		}

		return l
	}

	AddRoutes([]Route{
		{
			URI:    "test-public-endpoint",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: "test-public-endpoint"})
			},
			Prefix: PublicApi,
		},
		{
			URI:    "test-public-endpoint",
			Method: http.MethodPost,
			Prefix: PublicApi,
			Function: func(ctx WebContext) {
				body, err := ctx.StringBody()
				if err != nil {
					ctx.ErrorResponse(http.StatusInternalServerError, err)
					return
				}
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: body})
			},
		},
		{
			URI:    "test-private-endpoint/{msg}",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				msg := ctx.PathParam("msg")
				ctx.JsonResponse(http.StatusOK, &Resp{msg})
			},
			Prefix: PrivateApi,
		},
		{
			URI:    "test-authenticated-endpoint",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, ctx.AuthenticationContext())
			},
			Prefix: AuthenticatedApi,
		},
		{
			URI:    "test-authenticated-endpoint/{msg}",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				msg := ctx.PathParam("msg")
				ctx.JsonResponse(http.StatusOK, &Resp{msg})
			},
			Prefix: AuthenticatedApi,
		},
		{
			URI:    "test-no-prefix-endpoint",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: "test-no-prefix-endpoint"})
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-empty-body",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.EmptyResponse(http.StatusNoContent)
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-error-body",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.ErrorResponse(http.StatusInternalServerError, fmt.Errorf("test-error-body"))
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-request-header",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, ctx.RequestHeader("X-Id"))
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-request-headers",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, ctx.RequestHeaders())
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-query-param",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, ctx.QueryParam("size"))
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-query-array-param",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, ctx.QueryArrayParam("idList"))
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-decode-query-params",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				var result Query
				if err := ctx.DecodeQueryParams(&result); err != nil {
					ctx.ErrorResponse(http.StatusBadRequest, err)
				}
				ctx.JsonResponse(http.StatusOK, result)
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-decode-body",
			Method: http.MethodPost,
			Function: func(ctx WebContext) {
				var result Resp
				if err := ctx.DecodeBody(&result); err != nil {
					ctx.ErrorResponse(http.StatusBadRequest, err)
					return
				}
				ctx.JsonResponse(http.StatusOK, result)
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-context",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				msg := fmt.Sprintf("%s", ctx.Context())
				ctx.AddHeader("context", msg)
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: msg})
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "test-serve-file",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				path := test.MountAbsolutPath(test.REST_ENVIRONMENT_PATH) + "/resp.json"
				ctx.AddHeaders(map[string]string{"path": path})
				ctx.ServeFile(path)
			},
			Prefix: NoPrefix,
		},
		{
			URI:    "middleware/conflict",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: "conflict"})
			},
			Prefix: PublicApi,
		},
		{
			URI:    "middleware/success",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: "success"})
			},
			Prefix: PublicApi,
		},
		{
			URI:    "hook/conflict",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: "conflict"})
			},
			Prefix:      PublicApi,
			BeforeEnter: beforeEnterApply,
		},
		{
			URI:    "hook/success",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: "success"})
			},
			Prefix:      PublicApi,
			BeforeEnter: beforeEnterApply,
		},
	})

	l := listener()
	config.PORT = l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	Use(&customMiddlewareTest{})
	go ListenAndServe()
	time.Sleep(1 * time.Second)
	client := restclient.NewRestClient("test-server", fmt.Sprintf("http://localhost:%d", config.PORT), 1)

	t.Run("Should return status 200 (OK) in health-check", func(t *testing.T) {
		resp := restclient.Get[healtCheck](context.Background(), client, "/management/health", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.Equal(t, "OK", resp.Body().Status)
	})

	t.Run("Should return error 404 (not found) when endpoint not exists", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/not-exists-endpoint", nil)
		assert.NotNil(t, resp)
		assert.Error(t, resp.Error(), "404 statusCode")
		assert.Nil(t, resp.Body())
	})

	t.Run("Should return 200 (OK) in public api", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/public/test-public-endpoint", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, "test-public-endpoint", resp.Body().Msg)
	})

	t.Run("Should return 200 (OK) in private api", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/private/test-private-endpoint/abc", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, "abc", resp.Body().Msg)
	})

	t.Run("Should return error 401 (Unauthorized) when not inform the userId in authenticated api", func(t *testing.T) {
		r := restclient.Get[Resp](context.Background(), client, "/api/test-authenticated-endpoint/abc", map[string]string{tenantIDHeader: TenantId})
		assert.NotNil(t, r)
		assert.Error(t, r.Error(), "401 statusCode")
		assert.Nil(t, r.Body())
	})

	t.Run("Should return error 401 (Unauthorized) when not inform the tenantId in authenticated api", func(t *testing.T) {
		r := restclient.Get[Resp](context.Background(), client, "/api/test-authenticated-endpoint/abc", map[string]string{userIDHeader: UserId})
		assert.NotNil(t, r)
		assert.Error(t, r.Error(), "401 statusCode")
		assert.Nil(t, r.Body())
	})

	t.Run("Should return error 401 (Unauthorized) when not inform the credentials in authenticated api", func(t *testing.T) {
		r := restclient.Get[Resp](context.Background(), client, "/api/test-authenticated-endpoint/abc", nil)
		assert.NotNil(t, r)
		assert.Error(t, r.Error(), "401 statusCode")
		assert.Nil(t, r.Body())
	})

	t.Run("Should return 200 (OK) in authenticated api", func(t *testing.T) {
		r := restclient.Get[Resp](context.Background(), client, "/api/test-authenticated-endpoint/abc", map[string]string{tenantIDHeader: TenantId, userIDHeader: UserId})
		assert.NotNil(t, r)
		assert.NoError(t, r.Error())
		assert.NotNil(t, r.Body())
		assert.Equal(t, "abc", r.Body().Msg)
	})

	t.Run("Should return 200 (OK) in no prefix api", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/test-no-prefix-endpoint", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, "test-no-prefix-endpoint", resp.Body().Msg)
	})

	t.Run("Should return error response body", func(t *testing.T) {
		resp := restclient.Get[Error](context.Background(), client, "/test-error-body", nil)
		assert.NotNil(t, resp)
		assert.Error(t, resp.Error(), "500 statusCode")
		assert.Nil(t, resp.Body())
	})

	t.Run("Should return empty response body", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/test-empty-body", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.Nil(t, resp.Body())
	})

	t.Run("Should return request header", func(t *testing.T) {
		expected := []string{"123"}
		resp := restclient.Get[[]string](context.Background(), client, "/test-request-header", map[string]string{"X-Id": "123"})
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.Equal(t, &expected, resp.Body())
	})

	t.Run("Should return request headers", func(t *testing.T) {
		expected := []string{"456"}
		resp := restclient.Get[map[string][]string](context.Background(), client, "/test-request-headers", map[string]string{"X-Id": "456"})
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, expected, (*resp.Body())["X-Id"])
	})

	t.Run("Should return query param", func(t *testing.T) {
		expected := "10"
		resp := restclient.Get[string](context.Background(), client, "/test-query-param?size=10", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, expected, *resp.Body())
	})

	t.Run("Should return query array param", func(t *testing.T) {
		expected := []string{"10", "20", "30"}
		resp := restclient.Get[[]string](context.Background(), client, "/test-query-array-param?page=1&idList=10,20,30", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, expected, *resp.Body())
	})

	t.Run("Should return decoded query param", func(t *testing.T) {
		expected := &Query{Msg: "decoded", Size: 10}
		resp := restclient.Get[Query](context.Background(), client, "/test-decode-query-params?size=10&msg=decoded", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, expected, resp.Body())
	})

	t.Run("Should return bad request when an error occurred in decoded body", func(t *testing.T) {
		resp := restclient.Post[Resp](context.Background(), client, "/test-decode-body", &Resp{}, nil)
		assert.NotNil(t, resp)
		assert.Error(t, resp.Error())
		assert.Nil(t, resp.Body())
	})

	t.Run("Should return decoded body", func(t *testing.T) {
		expected := &Resp{Msg: "decoded"}
		resp := restclient.Post[Resp](context.Background(), client, "/test-decode-body", expected, nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, expected, resp.Body())
	})

	t.Run("Should return context", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/test-context", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.NotEmpty(t, resp.Body().Msg)
	})

	t.Run("Should return server file", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/test-serve-file", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, "test file", resp.Body().Msg)
	})

	t.Run("Should return error 200 using raw body", func(t *testing.T) {
		body := "text message"
		resp := restclient.Post[Resp](context.Background(), client, "/public/test-public-endpoint", &body, nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, "\"text message\"", resp.Body().Msg)
	})

	t.Run("should validate custom middleware with error", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/public/middleware/conflict", nil)
		assert.NotNil(t, resp)
		assert.Equal(t, "409 statusCode", resp.Error().Error())
		assert.Equal(t, http.StatusConflict, resp.Status())
	})

	t.Run("should validate custom middleware with success", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/public/middleware/success", nil)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.Status())
		assert.Equal(t, "success", resp.Body().Msg)
	})

	t.Run("should validate hook before enter route with error", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/public/hook/conflict", nil)
		assert.NotNil(t, resp)
		assert.Equal(t, "409 statusCode", resp.Error().Error())
		assert.Equal(t, http.StatusConflict, resp.Status())
	})

	t.Run("should validate hook before enter route with success", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/public/hook/success", nil)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.Status())
		assert.Equal(t, "success", resp.Body().Msg)
	})
}

func TestStartRestServerCustomAuthMiddleware(t *testing.T) {
	test.InitializeBaseTest()

	type Resp struct {
		Msg string `json:"msg" validate:"required"`
	}

	listener := func() (l net.Listener) {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Error("could not get available port")
		}
		return l
	}

	AddRoutes([]Route{
		{
			URI:    "users",
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.JsonResponse(http.StatusOK, &Resp{Msg: "test-custom-authentication-middleware"})
			},
			Prefix: AuthenticatedApi,
		},
	})

	l := listener()
	config.PORT = l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	CustomAuthMiddleware(&customAuthenticationContextMiddleware{})
	go ListenAndServe()
	time.Sleep(1 * time.Second)
	client := restclient.NewRestClient("test-server", fmt.Sprintf("http://localhost:%d", config.PORT), 1)

	t.Run("Should return status 200", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/api/users", map[string]string{"Authorization": "abcd1234"})
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, "test-custom-authentication-middleware", resp.Body().Msg)
	})

	t.Run("Should return error 401", func(t *testing.T) {
		resp := restclient.Get[Resp](context.Background(), client, "/api/users", nil)
		assert.NotNil(t, resp)
		assert.Error(t, resp.Error(), "401 statusCode")
		assert.Nil(t, resp.Body())
	})
}
