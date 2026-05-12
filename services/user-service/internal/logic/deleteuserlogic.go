package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/aicong/mine-dispatch/services/user-service/internal/model"
	"github.com/aicong/mine-dispatch/services/user-service/internal/svc"
)

type DeleteUserLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svc *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{ctx: ctx, svc: svc}
}

func (l *DeleteUserLogic) DeleteUser(in *userv1.DeleteUserRequest) (*userv1.UserResponse, error) {
	if err := l.svc.DB.Delete(&model.User{}, in.Id).Error; err != nil {
		return &userv1.UserResponse{Code: 500, Message: err.Error()}, nil
	}
	return &userv1.UserResponse{Code: 0, Message: "success"}, nil
}
