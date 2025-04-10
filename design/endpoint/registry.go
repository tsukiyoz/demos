package endpoint

import "context"

type Registry interface {
	Register(name string, addr string) error
	Unregister(name string) error
	DiscoverEndpoints(name string) ([]string, error)
	WatchEndpoints(ctx context.Context, name string) (chan []string, error)
}
