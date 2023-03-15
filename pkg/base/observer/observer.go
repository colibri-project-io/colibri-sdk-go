package observer

// Observer is the default graceful shutdown contract
type Observer interface {
	Close()
}
