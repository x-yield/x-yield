package plugins

type ListenerInterface interface {
	Put()
}

type Listener struct {
	ListenerInterface
}
