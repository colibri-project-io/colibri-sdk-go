package restclient

import "net/http"

// ResponseData interface to encapsulate response data with Status code, Headers, Success body, Error body and Error
type ResponseData[T any, E any] interface {
	StatusCode() int
	Headers() http.Header
	SuccessBody() *T
	ErrorBody() *E
	Error() error

	IsInformationalResponse() bool
	IsSuccessfulResponse() bool
	IsRedirectionMessage() bool
	IsClientErrorResponse() bool
	IsServerErrorResponse() bool
	HasError() bool
	HasSuccess() bool
}

// newResponseData creates a new ResponseData instance with the provided data.
//
// Parameters:
// - statusCode: the status code of the response
// - headers: the headers of the response
// - successBody: the success body of the response
// - errorBody: the error body of the response
// - err: the error associated with the response
//
// Returns a ResponseData instance.
func newResponseData[T any, E any](statusCode int, headers http.Header, successBody *T, errorBody *E, err error) ResponseData[T, E] {
	return &responseData[T, E]{
		statusCode:  statusCode,
		headers:     headers,
		successBody: successBody,
		errorBody:   errorBody,
		err:         err,
	}
}

// responseData implements the ResponseData interface.
type responseData[T any, E any] struct {
	statusCode  int
	headers     http.Header
	successBody *T
	errorBody   *E
	err         error
}

// StatusCode returns the status code of the response.
//
// No parameters.
// Returns an integer.
func (r *responseData[T, E]) StatusCode() int {
	return r.statusCode
}

// Headers returns the headers of the response.
//
// No parameters.
// Returns http.Header.
func (r *responseData[T, E]) Headers() http.Header {
	return r.headers
}

// SuccessBody returns the success body of the response.
//
// No parameters.
// Returns a pointer to the success body type.
func (r *responseData[T, E]) SuccessBody() *T {
	return r.successBody
}

// ErrorBody returns the error body of the response.
//
// No parameters.
// Returns a pointer to the error body.
func (r *responseData[T, E]) ErrorBody() *E {
	return r.errorBody
}

// Error returns the error associated with the response.
//
// No parameters.
// Returns an error.
func (r *responseData[T, E]) Error() error {
	return r.err
}

// IsInformationalResponse checks if the response status code is in the informational range (100-199).
//
// No parameters.
// Returns a boolean.
func (r *responseData[T, E]) IsInformationalResponse() bool {
	return r.statusCode >= 100 && r.statusCode <= 199
}

// IsSuccessfulResponse checks if the response status code is within the successful range (200-299).
//
// No parameters.
// Returns a boolean.
func (r *responseData[T, E]) IsSuccessfulResponse() bool {
	return r.statusCode >= 200 && r.statusCode <= 299
}

// IsRedirectionMessage checks if the response status code is in the redirection range (300-399).
//
// No parameters.
// Returns a boolean.
func (r *responseData[T, E]) IsRedirectionMessage() bool {
	return r.statusCode >= 300 && r.statusCode <= 399
}

// IsClientErrorResponse checks if the response status code is within the client error range (400-499).
//
// No parameters.
// Returns a boolean.
func (r *responseData[T, E]) IsClientErrorResponse() bool {
	return r.statusCode >= 400 && r.statusCode <= 499
}

// IsServerErrorResponse checks if the response status code is within the server error range (500-599).
//
// No parameters.
// Returns a boolean.
func (r *responseData[T, E]) IsServerErrorResponse() bool {
	return r.statusCode >= 500 && r.statusCode <= 599
}

// HasError checks if the response has an error based on the presence of err or errorBody.
// No parameters.
// Returns a boolean.
func (r *responseData[T, E]) HasError() bool {
	return r.err != nil || r.errorBody != nil
}

// HasSuccess checks if the response has no errors and no error body.
//
// No parameters.
// Returns a boolean.
func (r *responseData[T, E]) HasSuccess() bool {
	return r.err == nil && r.errorBody == nil
}
