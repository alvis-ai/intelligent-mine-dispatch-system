package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/aicong/mine-dispatch/services/user-service/internal/model"
	"github.com/aicong/mine-dispatch/services/user-service/internal/svc"
)

type GetUserLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svc *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{ctx: ctx, svc: svc}
}

func (l *GetUserLogic) GetUser(in *userv1.GetUserRequest) (*userv1.UserResponse, error) {
	var user model.User
	if err := l.svc.DB.First(&user, in.Id).Error; err != nil {
		return &userv1.UserResponse{Code: 404, Message: "user not found"}, nil
	}
	return &userv1.UserResponse{
		Code:    0,
		Message: "success",
		Data: &userv1.User{
			Id:       user.ID,
			Username: user.Username,
			RealName: user.RealName,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Status:   user.Status,
			MineId:   user.MineID,
		},
	}, nil
}
