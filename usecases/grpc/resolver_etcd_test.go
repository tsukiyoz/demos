package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"testing"
	"time"

	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	"go.etcd.io/etcd/client/v3/naming/endpoints"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

type EtcdTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (s *EtcdTestSuite) SetupSuite() {
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(s.T(), err)
	s.client = cli
}

func (s *EtcdTestSuite) TestServer() {
	lis, err := net.Listen("tcp", ":8090")
	require.NoError(s.T(), err)

	em, err := endpoints.NewManager(s.client, "service/user")
	require.NoError(s.T(), err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	kaCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	var ttl int64 = 6
	lease, err := s.client.Grant(kaCtx, ttl)
	require.NoError(s.T(), err)

	addr := "127.0.0.1:8090"
	key := fmt.Sprintf("service/user/%s", addr)
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(lease.ID))
	require.NoError(s.T(), err)

	// keep alive
	kaCtx, kaCancel := context.WithCancel(context.Background())
	go func() {
		aliveCh, err := s.client.KeepAlive(kaCtx, lease.ID)
		if err != nil {
			s.T().Log(err)
		}
		for alive := range aliveCh {
			s.T().Log(alive.String())
		}
	}()

	// update
	go func() {
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := em.AddEndpoint(ctx, key, endpoints.Endpoint{
				Addr:     addr,
				Metadata: now.String(),
			}, etcdv3.WithLease(lease.ID))
			if err != nil {
				s.T().Log(err)
			}
			cancel()
		}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, NewUserService("test"))

	err = server.Serve(lis)
	s.T().Log(err)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.T().Log("shutdown server ...")
	kaCancel()
	err = em.DeleteEndpoint(ctx, key)
	server.GracefulStop()
	s.client.Close()
	s.T().Log("server exited")
}

func (s *EtcdTestSuite) TestClient() {
	builder, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)

	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(builder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}

	client := NewUserServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
	require.NoError(s.T(), err)
	s.T().Log(resp.User)
}

func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}
