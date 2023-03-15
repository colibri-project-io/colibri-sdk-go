package restserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/validator"
	"github.com/gorilla/mux"
)

type GorillaWebContext struct {
	writer  http.ResponseWriter
	request *http.Request
}

func NewGorillaWebContext(writer http.ResponseWriter, request *http.Request) *GorillaWebContext {
	return &GorillaWebContext{writer: writer, request: request}
}

func (gCtx *GorillaWebContext) Context() context.Context {
	return gCtx.request.Context()
}

func (gCtx *GorillaWebContext) AuthenticationContext() *security.AuthenticationContext {
	return security.GetAuthenticationContext(gCtx.request.Context())
}

func (gCtx *GorillaWebContext) RequestHeader(key string) []string {
	return gCtx.request.Header[key]
}

func (gCtx *GorillaWebContext) RequestHeaders() (headers map[string][]string) {
	return gCtx.request.Header
}

func (gCtx *GorillaWebContext) PathParam(key string) string {
	return mux.Vars(gCtx.request)[key]
}

func (gCtx *GorillaWebContext) QueryParam(key string) string {
	return gCtx.request.URL.Query().Get(key)
}

func (gCtx *GorillaWebContext) QueryArrayParam(key string) []string {
	result := []string{}
	for _, value := range gCtx.request.URL.Query()[key] {
		result = append(result, strings.Split(value, ",")...)
	}

	return result
}

func (gCtx *GorillaWebContext) DecodeQueryParams(object any) error {
	if err := validator.FormDecode(object, gCtx.request.URL.Query()); err != nil {
		return err
	}

	return validator.Struct(object)
}

func (gCtx *GorillaWebContext) DecodeBody(object any) error {
	body, err := io.ReadAll(gCtx.request.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, object); err != nil {
		return err
	}

	return validator.Struct(object)
}

func (gCtx *GorillaWebContext) AddHeader(key string, value string) {
	gCtx.writer.Header().Add(key, value)
}

func (gCtx *GorillaWebContext) AddHeaders(headers map[string]string) {
	for key, value := range headers {
		gCtx.writer.Header().Add(key, value)
	}
}

func (gCtx *GorillaWebContext) Redirect(url string, statusCode int) {
	http.Redirect(gCtx.writer, gCtx.request, url, statusCode)
}

func (gCtx *GorillaWebContext) ServeFile(path string) {
	http.ServeFile(gCtx.writer, gCtx.request, path)
}

func (gCtx *GorillaWebContext) JsonResponse(statusCode int, body any) {
	gCtx.writer.Header().Add("Content-Type", "application/json")
	gCtx.writer.WriteHeader(statusCode)

	if err := json.NewEncoder(gCtx.writer).Encode(body); err != nil {
		logging.Error(err.Error())
	}
}

func (gCtx *GorillaWebContext) ErrorResponse(statusCode int, err error) {
	logging.Error("[%s] %s (%d): %v", gCtx.request.Method, gCtx.request.RequestURI, statusCode, err)
	gCtx.JsonResponse(statusCode, Error{err.Error()})
}

func (gCtx *GorillaWebContext) EmptyResponse(statusCode int) {
	gCtx.writer.WriteHeader(statusCode)
}

func (gCtx *GorillaWebContext) StringBody() (string, error) {
	rc := gCtx.request.Body
	b, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("could not read body: %w", err)
	}
	return string(b), err
}

func (gCtx *GorillaWebContext) Path() string {
	if gCtx.request.URL != nil {
		return gCtx.request.URL.Path
	}
	return ""
}

func (gCtx *GorillaWebContext) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return gCtx.request.FormFile(key)
}
