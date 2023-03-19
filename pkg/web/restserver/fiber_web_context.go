package restserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/validator"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"net/http"
	"strings"
)

type FiberWebContext struct {
	ctx         *fiber.Ctx
	ResponseErr error
}

func NewFiberWebContext(ctx *fiber.Ctx) *FiberWebContext {
	return &FiberWebContext{ctx: ctx}
}

func (f *FiberWebContext) Context() context.Context {
	return f.ctx.UserContext()
}

func (f *FiberWebContext) AuthenticationContext() *security.AuthenticationContext {
	return security.GetAuthenticationContext(f.Context())
}

func (f *FiberWebContext) RequestHeader(key string) []string {
	return []string{f.ctx.Get(key, "")}
}

func (f *FiberWebContext) RequestHeaders() map[string][]string {
	headers := make(map[string][]string)

	f.ctx.Context().Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = strings.Split(string(value), ";")
	})

	return headers
}

func (f *FiberWebContext) PathParam(key string) string {
	return f.ctx.Params(key)
}

func (f *FiberWebContext) QueryParam(key string) string {
	return f.ctx.Query(key)
}

func (f *FiberWebContext) QueryArrayParam(key string) []string {
	return strings.Split(f.ctx.Query(key), ",")
}

func (f *FiberWebContext) DecodeQueryParams(value any) error {
	if err := f.ctx.QueryParser(value); err != nil {
		return fmt.Errorf("could not decode query params: %w", err)
	}
	return validator.Struct(value)
}

func (f *FiberWebContext) DecodeBody(value any) error {
	body := f.ctx.Body()

	if err := json.Unmarshal(body, value); err != nil {
		return err
	}

	return validator.Struct(value)
}

func (f *FiberWebContext) AddHeader(key string, value string) {
	f.ctx.Response().Header.Add(key, value)
}

func (f *FiberWebContext) AddHeaders(headers map[string]string) {
	for key, value := range headers {
		f.ctx.Response().Header.Add(key, value)
	}
}

func (f *FiberWebContext) ServeFile(path string) {
	if err := f.ctx.SendFile(path); err != nil {
		f.ErrorResponse(http.StatusInternalServerError, err)
	}
}

func (f *FiberWebContext) JsonResponse(statusCode int, body any) {
	f.ctx.Response().SetStatusCode(statusCode)
	if err := f.ctx.JSON(body); err != nil {
		f.ErrorResponse(http.StatusInternalServerError, err)
	}
}

func (f *FiberWebContext) ErrorResponse(statusCode int, err error) {
	f.ResponseErr = err
	f.ctx.Response().SetStatusCode(statusCode)
	_ = f.ctx.JSON(Error{err.Error()})
}

func (f *FiberWebContext) EmptyResponse(statusCode int) {
	f.ctx.Response().SetStatusCode(statusCode)
}

func (f *FiberWebContext) IsError() bool {
	return f.ResponseErr != nil
}

func (f *FiberWebContext) Redirect(url string, statusCode int) {
	if err := f.ctx.Redirect(url, statusCode); err != nil {
		logging.Error("Could not set set redirect %s %d: %v", url, statusCode, err)
	}
}

func (f *FiberWebContext) StringBody() (string, error) {
	return string(f.ctx.Body()), nil
}

func (f *FiberWebContext) Path() string {
	return f.ctx.Path()
}

func (f *FiberWebContext) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	fileHeader, err := f.ctx.FormFile(key)

	if err != nil {
		return nil, nil, err
	}

	if file, err := fileHeader.Open(); err != nil {
		return nil, nil, err
	} else {
		return file, fileHeader, nil
	}
}
