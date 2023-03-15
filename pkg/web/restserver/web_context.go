package restserver

import (
	"context"
	"mime/multipart"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
)

// Error is the structure from default error response
type Error struct {
	Error string `json:"error"`
}

// WebContext is the contract to implements the web context in the webrest server
type WebContext interface {
	// Context returns the call context
	Context() context.Context
	// AuthenticationContext returns a pointer of authentication context if exists
	AuthenticationContext() *security.AuthenticationContext

	//RequestHeader returns a slice of http request header for key
	RequestHeader(key string) []string
	// RequestHeaders returns a slice of http request header by key
	RequestHeaders() map[string][]string
	// PathParam returns an url path param for key
	PathParam(key string) string
	// QueryParam returns a url query param for key
	QueryParam(key string) string
	// QueryArrayParam returns a url query array param for key
	QueryArrayParam(key string) []string
	// DecodeQueryParams converts all url query params into structured object
	DecodeQueryParams(object any) error
	// DecodeBody converts http request body into structured object
	DecodeBody(object any) error
	// StringBody return request raw body
	StringBody() (string, error)
	// Path return request path
	Path() string
	// FormFile returns the first file for the provided form key.
	FormFile(key string) (multipart.File, *multipart.FileHeader, error)

	// AddHeader add header into http response
	AddHeader(key string, value string)
	// AddHeaders add map of headers into http response
	AddHeaders(headers map[string]string)
	// Redirect http response with redirect headers
	Redirect(url string, statusCode int)
	// ServeFile serve a file into http response
	ServeFile(path string)
	// JsonResponse converts any object on json http response body
	JsonResponse(statusCode int, body any)
	// ErrorResponse converts any error on structured error http json response body
	ErrorResponse(statusCode int, err error)
	// EmptyResponse returns empty http json response body
	EmptyResponse(statusCode int)
}
