package webrest

import (
	"encoding/json"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/graceful-shutdown"
	"io"
	"log"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	PublicApi        RoutePrefix = "/public/"
	PrivateApi       RoutePrefix = "/private/"
	AuthenticatedApi RoutePrefix = "/api/"
)

type HealthCheck struct {
	Status string `json:"status"`
}

type RoutePrefix string

type Route struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
	Prefix   RoutePrefix
}

var (
	appRoutes       []Route
	srv             *http.Server
	bodyValidator   = validator.New()
	schemaValidator = schema.NewDecoder()
)

func AddRoutes(routes []Route) {
	appRoutes = append(appRoutes, routes...)
}

func create() *mux.Router {
	router := mux.NewRouter()

	router.Use(accessControlMiddleware)
	router.Use(authenticationContextMiddleware)

	for _, route := range appRoutes {
		routeUri := string(route.Prefix) + route.URI
		router.HandleFunc(monitoring.WrapHandleFunc(routeUri, route.Function)).Methods(route.Method)
	}

	return router
}

type restObserver struct{}

func (o restObserver) Close() {
	logging.Info("closing http server")
	if err := srv.Close(); err != nil {
		logging.Error("error when closing http server: %v", err)
	}
	srv = nil
}

func ListenAndServe() {
	addHealthCheckRoute()

	appRouter := create()

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.PORT),
		Handler: appRouter,
	}

	logging.Info("Service '%s' running in %d port", "WEB-REST", config.PORT)
	gracefulshutdown.Attach(restObserver{})
	log.Fatal(srv.ListenAndServe())
}

func addHealthCheckRoute() {
	health := &HealthCheck{"OK"}
	appRoutes = append(appRoutes, Route{
		URI:    "/health",
		Method: http.MethodGet,
		Function: func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(health)
		},
	})
}

func DecodeParams[T any](r *http.Request) (*T, error) {
	var result = new(T)
	err := schemaValidator.Decode(result, r.URL.Query())
	return result, err
}

func DecodeBody[T any](r *http.Request) (*T, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var result = new(T)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	if err := bodyValidator.Struct(result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetPathParam(r *http.Request, name string) string {
	return mux.Vars(r)[name]
}
