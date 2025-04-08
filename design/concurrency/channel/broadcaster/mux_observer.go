package broadcaster

const (
	register = iota
	unregister
	purge
)

type taggedReg[T any] struct {
	sub *subObserver[T]
	ch  chan<- T
	op  int
}

type MuxObserver[T any] struct {
	subs  map[*subObserver[T]]map[chan<- T]struct{}
	reg   chan taggedReg[T]
	input chan taggedObservation[T]
}

func NewMuxObserver[T any]() *MuxObserver[T] {
	m := &MuxObserver[T]{
		subs:  map[*subObserver[T]]map[chan<- T]struct{}{},
		reg:   make(chan taggedReg[T]),
		input: make(chan taggedObservation[T]),
	}
	go m.run()
	return m
}

func (m *MuxObserver[T]) doReg(tr taggedReg[T]) {
	mm, exists := m.subs[tr.sub]
	if !exists {
		mm = map[chan<- T]struct{}{}
		m.subs[tr.sub] = mm
	}
	mm[tr.ch] = struct{}{}
}

func (m *MuxObserver[T]) doUnreg(tr taggedReg[T]) {
	mm, exists := m.subs[tr.sub]
	if exists {
		delete(mm, tr.ch)
		if len(mm) == 0 {
			delete(m.subs, tr.sub)
		}
	}
}

func (m *MuxObserver[T]) handleReg(tr taggedReg[T]) {
	switch tr.op {
	case register:
		m.doReg(tr)
	case unregister:
		m.doUnreg(tr)
	case purge:
		delete(m.subs, tr.sub)
	}
}

func (m *MuxObserver[T]) run() {
	for {
		select {
		case tr, ok := <-m.reg:
			if ok {
				m.handleReg(tr)
			} else {
				return
			}
		default:
			select {
			case to := <-m.input:
				m.broadcast(to)
			case tr, ok := <-m.reg:
				if ok {
					m.handleReg(tr)
				} else {
					return
				}
			}
		}
	}
}

func (m *MuxObserver[T]) Close() {
	close(m.reg)
}

type taggedObservation[T any] struct {
	sub *subObserver[T]
	op  T
}

func (m *MuxObserver[T]) broadcast(to taggedObservation[T]) {
	for ch := range m.subs[to.sub] {
		ch <- to.op
	}
}

func (m *MuxObserver[T]) Sub() Broadcaster[T] {
	return &subObserver[T]{m}
}

type subObserver[T any] struct {
	mo *MuxObserver[T]
}

func (s *subObserver[T]) Register(ch chan<- T) {
	s.mo.reg <- taggedReg[T]{s, ch, register}
}

func (s *subObserver[T]) Unregister(ch chan<- T) {
	s.mo.reg <- taggedReg[T]{s, ch, unregister}
}

func (s *subObserver[T]) Close() error {
	s.mo.reg <- taggedReg[T]{s, nil, purge}
	return nil
}

func (s *subObserver[T]) Submit(ob T) {
	s.mo.input <- taggedObservation[T]{s, ob}
}

func (s *subObserver[T]) TrySubmit(ob T) bool {
	if s == nil {
		return false
	}
	select {
	case s.mo.input <- taggedObservation[T]{s, ob}:
		return true
	default:
		return false
	}
}
