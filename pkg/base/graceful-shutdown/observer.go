package gracefulshutdown

type Observer interface {
	Close()
}
