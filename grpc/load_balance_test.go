package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"testing"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/tsukaychan/demos/grpc/balancer/wrr"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/weightedroundrobin"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var prefix = "service/test"

type BalancerTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (b *BalancerTestSuite) SetupSuite() {
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(b.T(), err)
	b.client = cli
}

func (b *BalancerTestSuite) startWeightedServer(name, addr string, weight int) {
	// net listen
	lis, err := net.Listen("tcp", addr)
	require.NoError(b.T(), err)

	// args
	addr = b.getOutboundIP() + addr
	key := prefix + "/" + addr

	// endpoints manager
	em, err := endpoints.NewManager(b.client, prefix)
	require.NoError(b.T(), err)

	// lease
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var ttl int64 = 15
	lease, err := b.client.Grant(ctx, ttl)
	require.NoError(b.T(), err)

	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(lease.ID))
	require.NoError(b.T(), err)

	// keep alive
	kaCtx, kaCancel := context.WithCancel(context.Background())
	go func() {
		ch, err := b.client.KeepAlive(kaCtx, lease.ID)
		if err != nil {
			b.T().Log(err)
		}
		for msg := range ch {
			b.T().Log(msg.String())
		}
	}()

	// update
	go func() {
		ticker := time.NewTicker(time.Second)
		for tk := range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
				Addr: addr,
				Metadata: map[string]any{
					"time":   tk.String(),
					"weight": weight,
				},
			}, etcdv3.WithLease(lease.ID))
			if err != nil {
				b.T().Log(err)
			}
			cancel()
		}
	}()

	// setup server
	srv := grpc.NewServer()
	RegisterUserServiceServer(srv, NewUserService(name))

	// run server
	err = srv.Serve(lis)
	b.T().Log(err)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	// quit
	b.T().Log("shutdown server ...")
	kaCancel()

	// manager quit
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.DeleteEndpoint(ctx, key)

	// server quit
	srv.GracefulStop()

	// etcd client quit
	b.client.Close()
	b.T().Log("server exited")
}

func (b *BalancerTestSuite) getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func (b *BalancerTestSuite) TestPickFirst() {
	go func() {
		b.startWeightedServer("server1", ":8090", 10)
	}()
	go func() {
		b.startWeightedServer("server2", ":8091", 20)
	}()
	b.startWeightedServer("server3", ":8092", 30)
}

func (b *BalancerTestSuite) TestRClient() {
	builder, err := resolver.NewBuilder(b.client)
	require.NoError(b.T(), err)

	cc, err := grpc.Dial("etcd:///service/test",
		grpc.WithResolvers(builder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(b.T(), err)

	client := NewUserServiceClient(cc)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
		cancel()
		require.NoError(b.T(), err)
		b.T().Log(resp.User)
	}
}

func (b *BalancerTestSuite) TestRoundRobinClient() {
	builder, err := resolver.NewBuilder(b.client)
	require.NoError(b.T(), err)

	cc, err := grpc.Dial("etcd:///service/test",
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig("{\"loadBalancingConfig\":[{\"round_robin\":{}}]}"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(b.T(), err)

	client := NewUserServiceClient(cc)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
		cancel()
		if err != nil {
			b.T().Logf("get error: %v", err)
		} else {
			b.T().Log(resp.User)
		}
	}
}

// base on edf algo.
func (b *BalancerTestSuite) TestWeightedRoundRobinClient() {
	builder, err := resolver.NewBuilder(b.client)
	require.NoError(b.T(), err)

	cc, err := grpc.Dial("etcd:///service/test",
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig("{\"loadBalancingConfig\":[{\"weighted_round_robin\":{}}]}"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(b.T(), err)

	client := NewUserServiceClient(cc)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
		cancel()
		require.NoError(b.T(), err)
		b.T().Log(resp.User)
	}
}

func (b *BalancerTestSuite) TestCustomClient() {
	builder, err := resolver.NewBuilder(b.client)
	require.NoError(b.T(), err)

	cc, err := grpc.Dial("etcd:///service/test",
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig("{\"loadBalancingConfig\":[{\"custom\":{}}]}"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(b.T(), err)

	client := NewUserServiceClient(cc)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
		cancel()
		require.NoError(b.T(), err)
		b.T().Log(resp.User)
	}
}

func TestLoadBalance(t *testing.T) {
	suite.Run(t, new(BalancerTestSuite))
}
