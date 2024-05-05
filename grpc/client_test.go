package grpc

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestGrpcClient(t *testing.T) {
	gctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("user", "123"))
	cc, err := grpc.Dial(":8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := NewUserServiceClient(cc)

	ctx, cancel := context.WithTimeout(gctx, time.Second)
	defer cancel()
	resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
	require.NoError(t, err)
	t.Log(resp.User)
}
