package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcd_Watch(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 3 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchCh := client.Watch(ctx, "foo")
	res := make(map[string]struct{})
	listKeys := func() []string {
		keys := make([]string, 0, len(res))
		for k := range res {
			keys = append(keys, k)
		}
		return keys
	}
	_ = listKeys

	watchDone := make(chan struct{})

	go func() {
		defer close(watchDone)
		for evts := range watchCh {
			if evts.Canceled {
				return
			}

			for _, ev := range evts.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					t.Logf("Put: %s : %s", ev.Kv.Key, ev.Kv.Value)
					res[string(ev.Kv.Value)] = struct{}{}
					// t.Logf("Keys: %v", listKeys())
				case clientv3.EventTypeDelete:
					t.Logf("Delete: key: %s : %s", ev.Kv.Key, ev.Kv.Value)
					// delete(res, string(ev.Kv.Value))
					// t.Logf("Keys: %v", listKeys())
				}
			}
		}
		t.Log("Watch channel closed")
	}()

	var getResp *clientv3.GetResponse

	time.Sleep(1 * time.Second)
	_, err = client.Put(ctx, "foo", "bar")
	require.NoError(t, err)
	getResp, err = client.Get(ctx, "foo")
	require.NoError(t, err)
	t.Log("Get response:", getResp.Kvs)

	time.Sleep(1 * time.Second)
	_, err = client.Put(ctx, "foo", "baz")
	require.NoError(t, err)
	getResp, err = client.Get(ctx, "foo")
	require.NoError(t, err)
	t.Log("Get response:", getResp.Kvs)

	time.Sleep(1 * time.Second)
	_, err = client.Delete(ctx, "foo")
	require.NoError(t, err)
	getResp, err = client.Get(ctx, "foo")
	require.NoError(t, err)
	t.Log("Get response:", getResp.Kvs)

	time.Sleep(2 * time.Second)
	cancel()

	time.Sleep(2 * time.Second)
	<-watchDone
	t.Log("Watch done")
}
