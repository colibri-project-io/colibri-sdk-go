package restclient

import "net/http"

// ResponseData interface to encapsulate response data with Status, Body, Header and Error
type ResponseData[Body any] interface {
	Status() int
	Body() *Body
	Header() http.Header
	Error() error
}

func newResponseData[Body any](status int, body *Body, header http.Header, err error) ResponseData[Body] {
	return &responseData[Body]{status: status, body: body, header: header, err: err}
}

type responseData[B any] struct {
	status int
	body   *B
	header http.Header
	err    error
}

func (r *responseData[B]) Status() int {
	return r.status
}

func (r *responseData[B]) Body() *B {
	return r.body
}

func (r *responseData[B]) Header() http.Header {
	return r.header
}

func (r *responseData[B]) Error() error {
	return r.err
}
