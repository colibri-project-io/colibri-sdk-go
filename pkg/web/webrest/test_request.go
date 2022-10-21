package webrest

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	DefaultTestUserId   = "5e859dae-c879-11eb-b8bc-0242ac130003"
	DefaultTestTenantId = "5e859dae-c879-11eb-b8bc-0242ac130004"
)

type RestRequestTest struct {
	Method  string
	Url     string
	UrlVars map[string]string
	Headers map[string]string
	Body    io.Reader
}

func NewRestRequestTest(request *RestRequestTest) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(request.Method, request.Url, request.Body)
	req = mux.SetURLVars(req, request.UrlVars)

	req.Header.Add("X-TenantId", DefaultTestTenantId)
	req.Header.Add("X-UserId", DefaultTestUserId)
	req = req.WithContext(security.NewAuthenticationContext(uuid.MustParse(DefaultTestTenantId), uuid.MustParse(DefaultTestUserId)).SetInContext(req.Context()))

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	w := httptest.NewRecorder()
	return req, w
}
