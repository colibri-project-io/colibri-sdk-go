package restserver

import (
	"errors"
	"log"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/observer"
)

var (
	srvRoutes         []Route
	customMiddlewares []CustomMiddleware
	srv               Server
	customAuth        CustomAuthenticationMiddleware

	errUserUnauthenticated = errors.New("user not authenticated")
)

// Server is the contract to http server implementation
type Server interface {
	initialize()
	shutdown() error
	injectMiddlewares()
	injectCustomMiddlewares()
	injectRoutes()
	listenAndServe() error
}

// AddRoutes add list of routes in the webrest server
func AddRoutes(routes []Route) {
	srvRoutes = append(srvRoutes, routes...)
}

func CustomAuthMiddleware(fn CustomAuthenticationMiddleware) {
	customAuth = fn
}

func Use(m CustomMiddleware) {
	customMiddlewares = append(customMiddlewares, m)
}

// ListenAndServe initialize, configure and expose the web rest server
func ListenAndServe() {
	addHealthCheckRoute()
	addDocumentationRoute()

	srv = createFiberServer()
	srv.initialize()
	srv.injectMiddlewares()
	srv.injectCustomMiddlewares()
	srv.injectRoutes()

	observer.Attach(restObserver{})
	logging.Info("Service '%s' running in %d port", "WEB-REST", config.PORT)
	if err := srv.listenAndServe(); err != nil {
		log.Fatalf("Error rest server: %v", err)
	}
}
