package restserver

import "github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"

type restObserver struct {
}

func (o restObserver) Close() {
	logging.Info("closing http server")
	if err := srv.shutdown(); err != nil {
		logging.Error("error when closing http server: %v", err)
	}
	srv = nil
}
