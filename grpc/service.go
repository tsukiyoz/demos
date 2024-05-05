package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

var _ UserServiceServer = (*UserService)(nil)

type UserService struct {
	UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (svc *UserService) GetByID(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return &GetByIDResp{
			User: &User{
				Id:   req.Id,
				Name: "tsukiyo" + md.Get("user")[0],
			},
		}, nil
	}
	return &GetByIDResp{
		User: &User{
			Id:   req.Id,
			Name: "tsukiyo",
		},
	}, nil
}
