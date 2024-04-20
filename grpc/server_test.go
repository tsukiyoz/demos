package grpc

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
)

func TestUserServiceClient(t *testing.T) {
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	userServer := &UserService{}
	RegisterUserServiceServer(srv, userServer)
	listen, err := net.Listen("tcp", ":8090")
	require.NoError(t, err)
	err = srv.Serve(listen)
	require.NoError(t, err)
}
