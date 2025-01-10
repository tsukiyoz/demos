package epoll

import (
	"crypto/tls"
	"net"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type epoll struct {
	fd    int
	conns map[int]net.Conn
	lock  *sync.RWMutex
}

func NewEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		fd:    fd,
		lock:  &sync.RWMutex{},
		conns: make(map[int]net.Conn),
	}, nil
}

func (e *epoll) Add(conn net.Conn) error {
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.conns[fd] = conn
	return nil
}

func (e *epoll) Remove(conn net.Conn) error {
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.conns, fd)
	return nil
}

func (e *epoll) Wait() ([]net.Conn, error) {
	evts := make([]unix.EpollEvent, 100)
retry:
	n, err := unix.EpollWait(e.fd, evts, 100)
	if err != nil {
		if err == unix.EINTR {
			goto retry
		}
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	conns := make([]net.Conn, 0, n)
	for i := 0; i < n; i++ {
		conn := e.conns[int(evts[i].Fd)]
		conns = append(conns, conn)
	}
	return conns, nil
}

func socketFD(conn net.Conn) int {
	switch conn := conn.(type) {
	case syscall.Conn:
		rowConn, err := conn.SyscallConn()
		if err != nil {
			return 0
		}

		var fd int
		err = rowConn.Control(func(fileDescriptor uintptr) {
			fd = int(fileDescriptor)
		})
		if err != nil {
			return 0
		}

		return fd
	case *tls.Conn:
		tlsConn := conn.NetConn()
		return socketFD(tlsConn)
	default:
		panic("unsupported connection type")
	}
}
