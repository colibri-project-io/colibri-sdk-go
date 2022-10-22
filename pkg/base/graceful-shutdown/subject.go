package gracefulshutdown

import (
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"os"
	"os/signal"
	"syscall"
)

type subject interface {
	attach(observer Observer)
	notify()
}

var services subject

func Initialize() {
	ch := make(chan os.Signal, 1)
	services = &service{
		observers: make([]Observer, 0, 0),
	}
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL, os.Interrupt)

	go func() {
		sig := <-ch
		logging.Warn("notify shutdown: %+v", sig)
		services.notify()
	}()
}

func Attach(o Observer) {
	services.attach(o)
}

type service struct {
	observers []Observer
}

func (s *service) attach(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *service) notify() {
	for _, observer := range s.observers {
		observer.Close()
	}
}
