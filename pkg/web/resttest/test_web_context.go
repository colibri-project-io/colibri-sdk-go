package resttest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
	"github.com/go-playground/form/v4"
	playValidator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

var (
	validator   = playValidator.New()
	formDecoder = form.NewDecoder()
)

type testWebContext struct {
	req    *http.Request
	writer http.ResponseWriter
}

func (t *testWebContext) Context() context.Context {
	return t.req.Context()
}

func (t *testWebContext) AuthenticationContext() *security.AuthenticationContext {
	return security.GetAuthenticationContext(t.req.Context())
}

func (t *testWebContext) RequestHeader(key string) []string {
	return t.req.Header[key]
}

func (t *testWebContext) RequestHeaders() map[string][]string {
	return t.req.Header
}

func (t *testWebContext) PathParam(key string) string {
	return mux.Vars(t.req)[key]
}

func (t *testWebContext) QueryParam(key string) string {
	return t.req.URL.Query().Get(key)
}

func (t *testWebContext) QueryArrayParam(key string) []string {
	result := make([]string, 0)
	for _, value := range t.req.URL.Query()[key] {
		result = append(result, strings.Split(value, ",")...)
	}

	return result
}

func (t *testWebContext) DecodeQueryParams(object any) error {
	if err := formDecoder.Decode(object, t.req.URL.Query()); err != nil {
		return err
	}
	return validator.Struct(object)
}

func (t *testWebContext) DecodeBody(object any) error {
	body, err := io.ReadAll(t.req.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, object); err != nil {
		return err
	}

	return validator.Struct(object)
}

func (t *testWebContext) StringBody() (string, error) {
	rc := t.req.Body
	b, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("could not read body: %w", err)
	}
	return string(b), err
}

func (t *testWebContext) AddHeader(key string, value string) {
	t.writer.Header().Add(key, value)
}

func (t *testWebContext) AddHeaders(headers map[string]string) {
	for key, value := range headers {
		t.writer.Header().Add(key, value)
	}
}

func (t *testWebContext) Redirect(url string, statusCode int) {
	http.Redirect(t.writer, t.req, url, statusCode)
}

func (t *testWebContext) ServeFile(path string) {
	http.ServeFile(t.writer, t.req, path)
}

func (t *testWebContext) JsonResponse(statusCode int, body any) {
	t.writer.Header().Add("Content-Type", "application/json")
	t.writer.WriteHeader(statusCode)

	if err := json.NewEncoder(t.writer).Encode(body); err != nil {
		logging.Error(err.Error())
	}
}

func (t *testWebContext) ErrorResponse(statusCode int, err error) {
	logging.Error("[%s] %s (%d): %v", t.req.Method, t.req.RequestURI, statusCode, err)
	t.JsonResponse(statusCode, restserver.Error{Error: err.Error()})
}

func (t *testWebContext) EmptyResponse(statusCode int) {
	t.writer.WriteHeader(statusCode)
}

func (t *testWebContext) Path() string {
	return t.req.URL.Path
}

func (t *testWebContext) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return t.req.FormFile(key)
}
