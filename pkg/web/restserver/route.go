package restserver

import (
	"fmt"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"golang.org/x/exp/slices"
)

// RoutePrefix is the type from default's routes
type RoutePrefix string

const (
	PublicApi        RoutePrefix = "/public/"
	PrivateApi       RoutePrefix = "/private/"
	AuthenticatedApi RoutePrefix = "/api/"
	NoPrefix         RoutePrefix = "/"
)

// Route is the structure from inject the routes in the http router
type Route struct {
	URI         string
	Method      string
	Prefix      RoutePrefix
	Function    func(ctx WebContext)
	BeforeEnter func(ctx WebContext) *MiddlewareError
}

type healtCheck struct {
	Status string `json:"status"`
}

func addHealthCheckRoute() {
	const route = "/management/health"
	srvRoutes = append(srvRoutes, Route{
		URI:    route,
		Method: http.MethodGet,
		Function: func(ctx WebContext) {
			ctx.JsonResponse(http.StatusOK, &healtCheck{"OK"})
		},
	})

	logging.Info(fmt.Sprintf("Starting health-check on route: %s", route))
}

func addDocumentationRoute() {
	if slices.Contains([]string{config.ENVIRONMENT_SANDBOX, config.ENVIRONMENT_DEVELOPMENT}, config.ENVIRONMENT) {
		const route = "/v2/api-docs"
		srvRoutes = append(srvRoutes, Route{
			URI:    route,
			Method: http.MethodGet,
			Function: func(ctx WebContext) {
				ctx.ServeFile("./docs/swagger.json")
			},
		})

		logging.Info(fmt.Sprintf("Starting documentation on route: %s", route))
	}
}
