package server

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/aicong/mine-dispatch/services/user-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/user-service/internal/svc"
)

type UserServer struct {
	svc *svc.ServiceContext
	userv1.UnimplementedUserServiceServer
}

func NewUserServer(svc *svc.ServiceContext) *UserServer {
	return &UserServer{svc: svc}
}

func (s *UserServer) CreateUser(ctx context.Context, in *userv1.CreateUserRequest) (*userv1.UserResponse, error) {
	l := logic.NewCreateUserLogic(ctx, s.svc)
	return l.CreateUser(in)
}

func (s *UserServer) GetUser(ctx context.Context, in *userv1.GetUserRequest) (*userv1.UserResponse, error) {
	l := logic.NewGetUserLogic(ctx, s.svc)
	return l.GetUser(in)
}

func (s *UserServer) UpdateUser(ctx context.Context, in *userv1.UpdateUserRequest) (*userv1.UserResponse, error) {
	l := logic.NewUpdateUserLogic(ctx, s.svc)
	return l.UpdateUser(in)
}

func (s *UserServer) DeleteUser(ctx context.Context, in *userv1.DeleteUserRequest) (*userv1.UserResponse, error) {
	l := logic.NewDeleteUserLogic(ctx, s.svc)
	return l.DeleteUser(in)
}

func (s *UserServer) ListUser(ctx context.Context, in *userv1.ListUserRequest) (*userv1.UserListResponse, error) {
	l := logic.NewListUserLogic(ctx, s.svc)
	return l.ListUser(in)
}
