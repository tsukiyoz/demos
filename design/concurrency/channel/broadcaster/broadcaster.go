package broadcaster

type broadcaster[T any] struct {
	input   chan T
	reg     chan chan<- T
	unreg   chan chan<- T
	outputs map[chan<- T]struct{}
}

func NewBroadcaster[T any](buflen int) Broadcaster[T] {
	b := &broadcaster[T]{
		input:   make(chan T, buflen),
		reg:     make(chan chan<- T),
		unreg:   make(chan chan<- T),
		outputs: make(map[chan<- T]struct{}),
	}

	go b.run()

	return b
}

func (b *broadcaster[T]) run() {
	for {
		select {
		case m := <-b.input:
			b.broadcast(m)
		case ch, ok := <-b.reg:
			if ok {
				b.outputs[ch] = struct{}{}
			} else {
				return
			}
		case ch := <-b.unreg:
			delete(b.outputs, ch)
		}
	}
}

func (b *broadcaster[T]) broadcast(m T) {
	for ch := range b.outputs {
		ch <- m
	}
}

func (b *broadcaster[T]) Register(newch chan<- T) {
	b.reg <- newch
}

func (b *broadcaster[T]) Unregister(newch chan<- T) {
	b.unreg <- newch
}

func (b *broadcaster[T]) Close() error {
	close(b.reg)
	close(b.unreg)
	return nil
}

func (b *broadcaster[T]) Submit(m T) {
	if b != nil {
		b.input <- m
	}
}

func (b *broadcaster[T]) TrySubmit(m T) bool {
	if b == nil {
		return false
	}

	select {
	case b.input <- m:
		return true
	default:
		return false
	}
}
