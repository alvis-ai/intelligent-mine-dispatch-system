package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/pkg/utils"
	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/aicong/mine-dispatch/services/user-service/internal/model"
	"github.com/aicong/mine-dispatch/services/user-service/internal/svc"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svc *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{ctx: ctx, svc: svc}
}

func (l *CreateUserLogic) CreateUser(in *userv1.CreateUserRequest) (*userv1.UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return &userv1.UserResponse{Code: 500, Message: err.Error()}, nil
	}

	user := model.User{
		ID:       utils.NextID(),
		Username: in.Username,
		Password: string(hash),
		RealName: in.RealName,
		Email:    in.Email,
		Phone:    in.Phone,
		MineID:   in.MineId,
		Role:     1,
		Status:   1,
	}

	if err := l.svc.DB.Create(&user).Error; err != nil {
		return &userv1.UserResponse{Code: 500, Message: err.Error()}, nil
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
