package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/gorilla/mux"
)

type GorillaServer struct {
	engine *mux.Router
	srv    *http.Server
}

func createServer() *GorillaServer {
	return &GorillaServer{}
}

func (s *GorillaServer) initialize() {
	s.engine = mux.NewRouter()
}

func (s *GorillaServer) shutdown() error {
	return s.srv.Close()
}

func (s *GorillaServer) injectMiddlewares() {
	s.engine.Use(newRelicMiddleware)
	s.engine.Use(accessControlMiddleware)
	s.engine.Use(authenticationContextMiddleware)
}

func (s *GorillaServer) injectCustomMiddlewares() {
	for _, middleware := range customMiddlewares {
		s.registerCustomMiddleware(middleware)
	}
}

func (s *GorillaServer) injectRoutes() {
	for _, route := range srvRoutes {
		routeUri := string(route.Prefix) + route.URI
		fn := route.Function
		beforeEnter := route.BeforeEnter

		s.engine.HandleFunc(routeUri, func(w http.ResponseWriter, r *http.Request) {
			webCtx := &GorillaWebContext{writer: w, request: r}
			if beforeEnter != nil {
				if err := beforeEnter(webCtx); err != nil {
					webCtx.ErrorResponse(err.StatusCode, err.Err)
				}
			}
			fn(webCtx)
		}).Methods(route.Method)
	}
}

func (s *GorillaServer) listenAndServe() error {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.PORT),
		Handler: s.engine,
	}

	return s.srv.ListenAndServe()
}

func (s *GorillaServer) registerCustomMiddleware(middleware CustomMiddleware) {
	s.engine.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			webCtx := &GorillaWebContext{writer: w, request: r}
			if err := middleware.Apply(webCtx); err != nil {
				w.WriteHeader(err.StatusCode)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					logging.Error(err.Error())
				}
				return
			}
			next.ServeHTTP(w, r.WithContext(webCtx.Context()))
		})
	})
}
