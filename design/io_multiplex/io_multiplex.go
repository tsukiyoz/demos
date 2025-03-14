package iomultiplex

import "net"

type IOMultiplex interface {
	Add(conn net.Conn) error
	Remove(conn net.Conn) error
	Wait() ([]net.Conn, error)
}
