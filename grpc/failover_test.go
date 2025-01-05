package grpc

import (
	"context"
	_ "embed"
	"net"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	inet "github.com/tsukaychan/demos/net/ip"
	"go.etcd.io/etcd/client/v3/naming/endpoints"

	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/stretchr/testify/suite"
)

type FailoverTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (f *FailoverTestSuite) SetupSuite() {
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(f.T(), err)
	f.client = cli
}

//go:embed failover.json
var svsCfg string

func (f *FailoverTestSuite) TestServer() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		f.startServer(":8090", NewUserService("healthy"))
	}()

	go func() {
		defer wg.Done()

		f.startServer(":8091", NewFailService())
	}()

	wg.Wait()
}

func (f *FailoverTestSuite) TestClient() {
	builder, err := resolver.NewBuilder(f.client)
	require.NoError(f.T(), err)

	cc, err := grpc.Dial("etcd:///service/test",
		grpc.WithResolvers(builder),
		grpc.WithDefaultServiceConfig(svsCfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(f.T(), err)

	client := NewUserServiceClient(cc)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
		cancel()
		require.NoError(f.T(), err)
		f.T().Log(resp.User)
	}
}

func (f *FailoverTestSuite) startServer(addr string, svc UserServiceServer) {
	// net listen
	lis, err := net.Listen("tcp", addr)
	require.NoError(f.T(), err)

	// args
	addr = inet.GetOutboundIP() + addr
	key := prefix + "/" + addr

	// endpoints manager
	em, err := endpoints.NewManager(f.client, prefix)
	require.NoError(f.T(), err)

	// lease
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var ttl int64 = 15
	lease, err := f.client.Grant(ctx, ttl)
	require.NoError(f.T(), err)

	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(lease.ID))
	require.NoError(f.T(), err)

	// keep alive
	kaCtx, kaCancel := context.WithCancel(context.Background())
	go func() {
		ch, err := f.client.KeepAlive(kaCtx, lease.ID)
		if err != nil {
			f.T().Log(err)
		}
		for msg := range ch {
			f.T().Log(msg.String())
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
					"time": tk.String(),
				},
			}, etcdv3.WithLease(lease.ID))
			if err != nil {
				f.T().Log(err)
			}
			cancel()
		}
	}()

	// setup server
	srv := grpc.NewServer()
	RegisterUserServiceServer(srv, svc)
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())

	// run server
	err = srv.Serve(lis)
	f.T().Log(err)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	// quit
	f.T().Log("shutdown server ...")
	kaCancel()

	// manager quit
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.DeleteEndpoint(ctx, key)

	// server quit
	srv.GracefulStop()

	// etcd client quit
	f.client.Close()
	f.T().Log("server exited")
}

func TestFailover(t *testing.T) {
	suite.Run(t, new(FailoverTestSuite))
}
