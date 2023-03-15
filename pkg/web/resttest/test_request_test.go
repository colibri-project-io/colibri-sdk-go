package resttest

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

	context, writer := NewRequestTest(&RequestTest{Method: http.MethodGet, Url: "http://teste.com/url", Headers: map[string]string{headerKey: headerValue}})
	assert.Equal(t, []string{headerValue}, context.RequestHeader(headerKey))
	assert.NotNil(t, writer)

	authContext := context.AuthenticationContext()
	assert.NotNil(t, authContext)
	assert.Equal(t, authContext.GetTenantID().String(), DefaultTestTenantId)
	assert.Equal(t, authContext.GetUserID().String(), DefaultTestUserId)
}
