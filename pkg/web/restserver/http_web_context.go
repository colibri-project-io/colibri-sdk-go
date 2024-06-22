package restserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/validator"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type httpWebContext struct {
	w   http.ResponseWriter
	req *http.Request
	err error
}

func newHttpWebContext(w http.ResponseWriter, req *http.Request) WebContext {
	return &httpWebContext{w: w, req: req}
}

func (ctx *httpWebContext) Context() context.Context {
	return ctx.req.Context()
}

func (ctx *httpWebContext) AuthenticationContext() *security.AuthenticationContext {
	return security.GetAuthenticationContext(ctx.req.Context())
}

func (ctx *httpWebContext) RequestHeader(key string) []string {
	return ctx.req.Header[key]
}

func (ctx *httpWebContext) RequestHeaders() map[string][]string {
	return ctx.req.Header
}

func (ctx *httpWebContext) PathParam(key string) string {
	return ctx.req.PathValue(key)
}

func (ctx *httpWebContext) QueryParam(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *httpWebContext) QueryArrayParam(key string) []string {
	result := []string{}
	for _, value := range ctx.req.URL.Query()[key] {
		result = append(result, strings.Split(value, ",")...)
	}

	return result
}

func (ctx *httpWebContext) DecodeQueryParams(object any) error {
	if err := validator.FormDecode(object, ctx.req.URL.Query()); err != nil {
		return err
	}

	return validator.Struct(object)
}

func (ctx *httpWebContext) DecodeBody(object any) error {
	body, err := io.ReadAll(ctx.req.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, object); err != nil {
		return err
	}

	return validator.Struct(object)
}

func (ctx *httpWebContext) StringBody() (string, error) {
	rc := ctx.req.Body
	b, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("could not read body: %w", err)
	}
	return string(b), err
}

func (ctx *httpWebContext) Path() string {
	if ctx.req.URL != nil {
		return ctx.req.URL.Path
	}
	return ""
}

func (ctx *httpWebContext) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.req.FormFile(key)
}

func (ctx *httpWebContext) AddHeader(key string, value string) {
	ctx.w.Header().Add(key, value)
}

func (ctx *httpWebContext) AddHeaders(headers map[string]string) {
	for key, value := range headers {
		ctx.w.Header().Add(key, value)
	}
}

func (ctx *httpWebContext) Redirect(url string, statusCode int) {
	http.Redirect(ctx.w, ctx.req, url, statusCode)
}

func (ctx *httpWebContext) ServeFile(path string) {
	http.ServeFile(ctx.w, ctx.req, path)
}

func (ctx *httpWebContext) JsonResponse(statusCode int, body any) {
	ctx.w.Header().Add("Content-Type", "application/json")
	ctx.w.WriteHeader(statusCode)

	if err := json.NewEncoder(ctx.w).Encode(body); err != nil {
		logging.Error(err.Error())
	}
}

func (ctx *httpWebContext) ErrorResponse(statusCode int, err error) {
	logging.Error("[%s] %s (%d): %v", ctx.req.Method, ctx.req.RequestURI, statusCode, err)
	ctx.JsonResponse(statusCode, Error{err.Error()})
}

func (ctx *httpWebContext) EmptyResponse(statusCode int) {
	ctx.w.WriteHeader(statusCode)
}
