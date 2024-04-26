package restserver

import "net/http"

type httpWebServer struct {
	srv *http.Server
}

func newHttpWebServer() Server {
	return &httpWebServer{}
}

func (srv *httpWebServer) initialize() {
	//TODO implement me
	panic("implement me")
}

func (srv *httpWebServer) shutdown() error {
	//TODO implement me
	panic("implement me")
}

func (srv *httpWebServer) injectMiddlewares() {
	//TODO implement me
	panic("implement me")
}

func (srv *httpWebServer) injectCustomMiddlewares() {
	//TODO implement me
	panic("implement me")
}

func (srv *httpWebServer) injectRoutes() {
	//TODO implement me
	panic("implement me")
}

func (srv *httpWebServer) listenAndServe() error {
	//TODO implement me
	panic("implement me")
}
