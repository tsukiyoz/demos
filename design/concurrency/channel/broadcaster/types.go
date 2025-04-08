package broadcaster

type Broadcaster[T any] interface {
	Register(chan<- T)
	Unregister(chan<- T)
	Close() error
	Submit(T)
	TrySubmit(T) bool
}
