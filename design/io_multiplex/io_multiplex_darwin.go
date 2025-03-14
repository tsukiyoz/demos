//go:build darwin

package iomultiplex

import (
	"crypto/tls"
	"net"
	"sync"
	"syscall"
)

type kqueue struct {
	fd    int
	conns map[int]net.Conn
	lock  *sync.RWMutex
}

func NewKqueue() (*kqueue, error) {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}
	return &kqueue{
		fd:    fd,
		lock:  &sync.RWMutex{},
		conns: make(map[int]net.Conn),
	}, nil
}

func (k *kqueue) Add(conn net.Conn) error {
	fd := socketFD(conn)
	ev := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ, // 监听可读事件
		Flags:  syscall.EV_ADD,      // 添加事件
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	// 将事件添加到 kqueue
	changeEvent := []syscall.Kevent_t{ev}
	_, err := syscall.Kevent(k.fd, changeEvent, nil, nil)
	if err != nil {
		return err
	}

	k.lock.Lock()
	defer k.lock.Unlock()
	k.conns[fd] = conn
	return nil
}

func (k *kqueue) Remove(conn net.Conn) error {
	fd := socketFD(conn)
	ev := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_DELETE, // 删除事件
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	// 从 kqueue 中删除事件
	changeEvent := []syscall.Kevent_t{ev}
	_, err := syscall.Kevent(k.fd, changeEvent, nil, nil)
	if err != nil {
		return err
	}

	k.lock.Lock()
	defer k.lock.Unlock()
	delete(k.conns, fd)
	return nil
}

func (k *kqueue) Wait() ([]net.Conn, error) {
	events := make([]syscall.Kevent_t, 100)
retry:
	n, err := syscall.Kevent(k.fd, nil, events, nil)
	if err != nil {
		if err == syscall.EINTR {
			goto retry
		}
		return nil, err
	}

	k.lock.RLock()
	defer k.lock.RUnlock()
	conns := make([]net.Conn, 0, n)
	for i := 0; i < n; i++ {
		fd := int(events[i].Ident)
		if conn, ok := k.conns[fd]; ok {
			conns = append(conns, conn)
		}
	}
	return conns, nil
}

func socketFD(conn net.Conn) int {
	switch conn := conn.(type) {
	case syscall.Conn:
		rawConn, err := conn.SyscallConn()
		if err != nil {
			return 0
		}

		var fd int
		err = rawConn.Control(func(fileDescriptor uintptr) {
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
