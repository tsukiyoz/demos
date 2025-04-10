package endpoint

import (
	"context"

	etcdv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRegistry struct {
	client *etcdv3.Client
}

// DiscoverEndpoints implements Registry.
func (e *EtcdRegistry) DiscoverEndpoints(name string) ([]string, error) {
	panic("unimplemented")
}

// Register implements Registry.
func (e *EtcdRegistry) Register(name string, addr string) error {
	panic("unimplemented")
}

// Unregister implements Registry.
func (e *EtcdRegistry) Unregister(name string) error {
	panic("unimplemented")
}

// WatchEndpoints implements Registry.
func (e *EtcdRegistry) WatchEndpoints(ctx context.Context, name string) (chan []string, error) {
	panic("unimplemented")
}

var _ Registry = (*EtcdRegistry)(nil)
