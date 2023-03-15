package resttest

import (
	"bytes"
	"net/http/httptest"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	DefaultTestUserId   = "5e859dae-c879-11eb-b8bc-0242ac130003"
	DefaultTestTenantId = "5e859dae-c879-11eb-b8bc-0242ac130004"
)

// RequestTest is a contract to test http requests
type RequestTest struct {
	Method  string
	Url     string
	UrlVars map[string]string
	Headers map[string]string
	Body    string
}

// NewRequestTest returns a new web context and response recorder from http tests
func NewRequestTest(request *RequestTest) (ctx restserver.WebContext, response *RequestTestResponse) {
	req := httptest.NewRequest(request.Method, request.Url, bytes.NewBuffer([]byte(request.Body)))
	req = mux.SetURLVars(req, request.UrlVars)

	req.Header.Add("X-TenantId", DefaultTestTenantId)
	req.Header.Add("X-UserId", DefaultTestUserId)
	req = req.WithContext(security.
		NewAuthenticationContext(uuid.MustParse(DefaultTestTenantId), uuid.MustParse(DefaultTestUserId)).
		SetInContext(req.Context()))

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	response = &RequestTestResponse{httptest.NewRecorder()}
	ctx = restserver.NewGorillaWebContext(response.writer, req)
	return
}
