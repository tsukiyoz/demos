package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestGrpcClient(t *testing.T) {
	// init
	cc, err := grpc.Dial(":8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authClientInterceptor),
	)
	require.NoError(t, err)
	client := NewUserServiceClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// mock get user id from jwt
	userID := uuid.New().String()
	ctx = context.WithValue(ctx, vkey{}, userID)

	resp, err := client.GetByID(ctx, &GetByIDReq{Id: 123})
	require.NoError(t, err)
	t.Log(resp.User)
	t.Log(resp.Msg)
}

var authClientInterceptor grpc.UnaryClientInterceptor = func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	uid, _ := ctx.Value("uid").(string)
	gctx := metadata.NewOutgoingContext(ctx, metadata.Pairs("uid", uid))

	return invoker(gctx, method, req, reply, cc, opts...)
}
