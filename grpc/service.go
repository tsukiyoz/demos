package grpc

import "context"

var _ UserServiceServer = (*UserService)(nil)

type UserService struct {
	UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (svc *UserService) GetByID(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error) {
	return &GetByIDResp{
		User: &User{
			Id:   req.Id,
			Name: "tsukiyo",
		},
	}, nil
}
