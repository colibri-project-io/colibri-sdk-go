package restserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestTest(t *testing.T) {
	const (
		headerKey   string = "Header-Test-Key"
		headerValue string = "Header-Test-Value"
	)

	handlerFn := func(ctx WebContext) {
		ctx.EmptyResponse(http.StatusNoContent)
	}

	reqTest := &RequestTest{Method: http.MethodGet, Url: "/url", Headers: map[string]string{headerKey: headerValue}}
	resp := NewRequestTest(reqTest, handlerFn)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}
