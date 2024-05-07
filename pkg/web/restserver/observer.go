package restserver

import (
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/observer"
)

type restObserver struct {
}

func (o restObserver) Close() {
	logging.Info("waiting to safely close the http server")
	if observer.WaitRunningTimeout() {
		logging.Warn("WaitGroup timed out, forcing close http server")
	}

	logging.Info("closing http server")
	if err := srv.shutdown(); err != nil {
		logging.Error("error when closing http server: %v", err)
	}

	srv = nil
}
