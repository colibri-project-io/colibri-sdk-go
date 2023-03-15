package observer

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type subject interface {
	attach(observer Observer)
	notify()
}

var services subject

// Initialize starts the subject observabilty
func Initialize() {
	ch := make(chan os.Signal, 1)
	services = &service{
		observers: make([]Observer, 0),
	}
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	go func() {
		sig := <-ch
		logging.Warn("notify shutdown: %+v", sig)
		services.notify()
	}()
}

// Attach attach the subject on services observer
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
