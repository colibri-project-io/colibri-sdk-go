package resttest

import (
	"bytes"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/web/restserver"
	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/google/uuid"
)

const (
	DefaultTestUserId   = "5e859dae-c879-11eb-b8bc-0242ac130003"
	DefaultTestTenantId = "5e859dae-c879-11eb-b8bc-0242ac130004"
)

// RequestTest is a contract to test http requests
type RequestTest struct {
	Method string
	Url    string
	//UrlVars map[string]string
	Headers map[string]string
	Body    string
}

// NewRequestTest returns a TestResponse with result of test execution
func NewRequestTest(request *RequestTest, handlerFn func(ctx restserver.WebContext)) *TestResponse {
	app := fiber.New()
	routeUri := convertUriToFiberUri(request.Url)
	app.Add(request.Method, routeUri, func(ctx *fiber.Ctx) error {
		webContext := restserver.NewFiberWebContext(ctx)

		handlerFn(webContext)

		if webContext.IsError() {
			return webContext.ResponseErr
		}
		return nil
	})

	req := httptest.NewRequest(request.Method, request.Url, bytes.NewBuffer([]byte(request.Body)))

	req.Header.Add("X-TenantId", DefaultTestTenantId)
	req.Header.Add("X-UserId", DefaultTestUserId)
	req = req.WithContext(security.NewAuthenticationContext(uuid.MustParse(DefaultTestTenantId), uuid.MustParse(DefaultTestUserId)).
		SetInContext(req.Context()))

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		panic(err)
	}

	return &TestResponse{resp: resp}
}

func convertUriToFiberUri(uri string) string {

	replacer := strings.NewReplacer(
		"{", "",
		"}", "",
	)

	paths := strings.Split(uri, "/")

	for idx, path := range paths {
		if pathIsPathParam(path) {
			paths[idx] = fmt.Sprintf(":%s", replacer.Replace(path))
		}
	}

	return strings.Join(paths, "/")
}

func pathIsPathParam(path string) bool {
	return strings.Contains(path, "{")
}
