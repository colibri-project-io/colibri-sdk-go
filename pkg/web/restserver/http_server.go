package restserver

import (
	"context"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"net/http"
)

type httpWebServer struct {
	srv *http.Server
}

func newHttpWebServer() Server {
	return &httpWebServer{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.PORT),
			Handler: http.NewServeMux(),
		},
	}
}

func (h *httpWebServer) initialize() {
	logging.Debug("Initializing HTTP Server at port ")
}

func (h *httpWebServer) shutdown() error {
	return h.srv.Shutdown(context.Background())
}

func (h *httpWebServer) injectMiddlewares() {
	//TODO implement me
	panic("implement me")
}

func (h *httpWebServer) injectCustomMiddlewares() {
	//TODO implement me
	panic("implement me")
}

func (h *httpWebServer) injectRoutes() {
	//TODO implement me
	panic("implement me")
}

func (h *httpWebServer) listenAndServe() error {
	return h.srv.ListenAndServe()
}
