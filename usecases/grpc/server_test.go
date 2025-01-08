package grpc

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
)

func TestUserServiceClient(t *testing.T) {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authServerInterceptor),
	)
	defer srv.GracefulStop()
	userServer := &UserService{
		Name: "health service",
	}
	RegisterUserServiceServer(srv, userServer)
	listen, err := net.Listen("tcp", ":8090")
	require.NoError(t, err)
	err = srv.Serve(listen)
	require.NoError(t, err)
}

var authServerInterceptor grpc.UnaryServerInterceptor = func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		uid := md.Get("uid")
		if len(uid) != 0 {
			ctx = context.WithValue(ctx, "uid", uid[0])
		}
	}
	return handler(ctx, req)
}
