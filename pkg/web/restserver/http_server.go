package restserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"strings"
)

type httpWebServer struct {
	srv     *http.Server
	mux     *http.ServeMux
	handler http.Handler
}

func newHttpWebServer() Server {
	return &httpWebServer{
		mux: http.NewServeMux(),
	}
}

func (h *httpWebServer) initialize() {
	logging.Debug("Initializing HTTP Server at port %d", config.PORT)
}

func (h *httpWebServer) shutdown() error {
	if h.srv == nil {
		panic("Server is nil")
	}
	logging.Debug("Shutting down HTTP Server")
	return h.srv.Shutdown(context.Background())
}

func (h *httpWebServer) injectMiddlewares() {
	h.handler = h.corsMiddleware(h.newRelicMiddleware(h.mux))

	if customAuth != nil {
		h.customAuthenticationMiddleware()
	} else {
		h.handler = h.authenticationContextMiddleware(h.handler)
	}
}

func (h *httpWebServer) injectCustomMiddlewares() {
	for _, middleware := range customMiddlewares {
		h.registerCustomMiddleware(middleware)
	}
}

func (h *httpWebServer) injectRoutes() {
	h.addMetricsRoute()
	h.addSwaggerUI()

	for _, route := range srvRoutes {
		routeUri := string(route.Prefix) + route.URI
		fn := route.Function
		beforeEnter := route.BeforeEnter

		h.mux.HandleFunc(fmt.Sprintf("%s %s", strings.ToUpper(route.Method), routeUri), func(w http.ResponseWriter, r *http.Request) {
			webContext := newHttpWebContext(w, r)
			if beforeEnter != nil {
				if err := beforeEnter(webContext); err != nil {
					w.WriteHeader(err.StatusCode)
					_ = json.NewEncoder(w).Encode(err)
					return
				}
				fn(webContext)
				return
			} else {
				fn(webContext)
			}
		})
	}
}

func (h *httpWebServer) listenAndServe() error {
	h.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.PORT),
		Handler: h.handler,
	}

	defer func() {
		if p := recover(); p != nil {
			logging.Error("panic recovering: %v", p)
		}
	}()

	return h.srv.ListenAndServe()
}

func (h *httpWebServer) addMetricsRoute() {
	const route = "/metrics"
	h.mux.Handle(fmt.Sprintf("GET %s", route), promhttp.Handler())
	logging.Debug(fmt.Sprintf("Starting metrics on route: %s", route))
}

func (h *httpWebServer) addSwaggerUI() {
	if config.IsDevelopmentEnvironment() {
		h.mux.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)
	}
}

func (h *httpWebServer) newRelicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("X-Request-URI", r.RequestURI)
		r.Header.Add("X-Protocol", r.Proto)
		txn, ctx := monitoring.StartWebRequest(r.Context(), r.Header, r.URL.Path, r.Method)
		defer monitoring.EndTransaction(txn)

		r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (h *httpWebServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *httpWebServer) registerCustomMiddleware(middleware CustomMiddleware) {
	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			webCtx := newHttpWebContext(w, r)
			if err := middleware.Apply(webCtx); err != nil {
				w.WriteHeader(err.StatusCode)
				_ = json.NewEncoder(w).Encode(err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	h.handler = fn(h.handler)
}

func (h *httpWebServer) customAuthenticationMiddleware() {
	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			webCtx := newHttpWebContext(w, r)
			if !strings.Contains(webCtx.Path(), string(AuthenticatedApi)) {
				next.ServeHTTP(w, r)
				return
			}

			authCtx, errAuth := customAuth.Apply(webCtx)
			if errAuth != nil {
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(errAuth); err != nil {
					logging.Error(err.Error())
				}
				return
			}

			if authCtx != nil && authCtx.Valid() {
				next.ServeHTTP(w, r.WithContext(authCtx.SetInContext(r.Context())))
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(UserUnauthenticatedError); err != nil {
				logging.Error(err.Error())
			}
		})
	}

	h.handler = fn(h.handler)
}

func (h *httpWebServer) authenticationContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, string(AuthenticatedApi)) {
			next.ServeHTTP(w, r)
			return
		}

		tenantID := r.Header.Get(tenantIDHeader)
		userID := r.Header.Get(userIDHeader)
		authCtx := security.NewAuthenticationContext(tenantID, userID)
		if authCtx.Valid() {
			next.ServeHTTP(w, r.WithContext(authCtx.SetInContext(r.Context())))
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(UserUnauthenticatedError); err != nil {
			logging.Error(err.Error())
		}
	})
}
