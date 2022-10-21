package webrest

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestRestServer(t *testing.T) {
	type Response struct {
		Msg string `json:"msg"`
	}

	listener := func() (l net.Listener) {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Error("could not get available port")
		}

		return l
	}

	monitoring.Initialize()
	AddRoutes([]Route{
		{
			URI:    "test-public-endpoint/",
			Method: "GET",
			Function: func(writer http.ResponseWriter, request *http.Request) {
				JsonResponse(writer, http.StatusOK, &Response{Msg: "test-public-endpoint"})
			},
			Prefix: PublicApi,
		},
		{
			URI:    "test-private-endpoint/{msg}",
			Method: "GET",
			Function: func(writer http.ResponseWriter, request *http.Request) {
				msg := GetPathParam(request, "msg")
				JsonResponse(writer, http.StatusOK, &Response{msg})
			},
			Prefix: PrivateApi,
		},
		{
			URI:    "test-authenticated-endpoint/{msg}",
			Method: "GET",
			Function: func(writer http.ResponseWriter, request *http.Request) {
				msg := GetPathParam(request, "msg")
				JsonResponse(writer, http.StatusOK, &Response{msg})
			},
			Prefix: AuthenticatedApi,
		},
	})

	l := listener()
	config.PORT = l.Addr().(*net.TCPAddr).Port
	l.Close()

	go ListenAndServe()
	time.Sleep(1 * time.Second)
	client := NewRestClient("test-server", fmt.Sprintf("http://localhost:%d", config.PORT), 1)

	t.Run("Should return status 200 (OK) in health-check", func(t *testing.T) {
		resp, err := Get[HealtCheck](context.Background(), client, "/health", nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "OK", resp.Status)
	})

	t.Run("Should return error 404 (not found) when endpoint not exists", func(t *testing.T) {
		resp, err := Get[Response](context.Background(), client, "/not-exists-endpoint", nil)
		assert.Error(t, err, "404 statusCode")
		assert.Nil(t, resp)
	})

	t.Run("Should return 200 (OK) in private api", func(t *testing.T) {
		resp, err := Get[Response](context.Background(), client, "/private/test-private-endpoint/abc", nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "abc", resp.Msg)
	})

	t.Run("Should return error 401 (Unauthorized) when not inform the userId in authenticated api", func(t *testing.T) {
		r, err := Get[Response](context.Background(), client, "/api/test-authenticated-endpoint/abc", map[string]string{"X-TenantId": DefaultTestTenantId})
		assert.Error(t, err, "401 statusCode")
		assert.Nil(t, r)
	})

	t.Run("Should return error 401 (Unauthorized) when not inform the tenantId in authenticated api", func(t *testing.T) {
		r, err := Get[Response](context.Background(), client, "/api/test-authenticated-endpoint/abc", map[string]string{"X-UserId": DefaultTestUserId})
		assert.Error(t, err, "401 statusCode")
		assert.Nil(t, r)
	})

	t.Run("Should return error 401 (Unauthorized) when not inform the credentials in authenticated api", func(t *testing.T) {
		r, err := Get[Response](context.Background(), client, "/api/test-authenticated-endpoint/abc", nil)
		assert.Error(t, err, "401 statusCode")
		assert.Nil(t, r)
	})

	t.Run("Should return 200 (OK) in authenticated api", func(t *testing.T) {
		r, err := Get[Response](context.Background(), client, "/api/test-authenticated-endpoint/abc", map[string]string{"X-TenantId": DefaultTestTenantId, "X-UserId": DefaultTestUserId})
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "abc", r.Msg)
	})
}

func TestDecodeParams(t *testing.T) {
	type RequestQueryParams struct {
		Page uint16 `schema:"page"`
		Size uint16 `schema:"size"`
	}

	r, _ := http.NewRequest(http.MethodGet, "/my-endpoint?page=1&size=10", nil)
	result, err := DecodeParams[RequestQueryParams](r)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint16(1), result.Page)
	assert.Equal(t, uint16(10), result.Size)
}

func TestDecodeBody(t *testing.T) {
	type RequestDecodeBody struct {
		Name string `json:"name"`
		Age  uint16 `json:"age"`
	}

	expected := &RequestDecodeBody{Name: "Zezinho", Age: 25}
	json, _ := json.Marshal(expected)

	t.Run("Should decode body with success", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/my-endpoint", strings.NewReader(string(json)))
		result, err := DecodeBody[RequestDecodeBody](r)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected, result)
	})

	t.Run("Should return error when body is a invalid json", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/my-endpoint", strings.NewReader("invalid json body"))
		result, err := DecodeBody[RequestDecodeBody](r)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
