package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/aicong/mine-dispatch/services/user-service/internal/model"
	"github.com/aicong/mine-dispatch/services/user-service/internal/svc"
)

type UpdateUserLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svc *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{ctx: ctx, svc: svc}
}

func (l *UpdateUserLogic) UpdateUser(in *userv1.UpdateUserRequest) (*userv1.UserResponse, error) {
	updates := map[string]interface{}{}
	if in.RealName != "" {
		updates["real_name"] = in.RealName
	}
	if in.Email != "" {
		updates["email"] = in.Email
	}
	if in.Phone != "" {
		updates["phone"] = in.Phone
	}
	if in.Role != 0 {
		updates["role"] = in.Role
	}
	if in.Status != 0 {
		updates["status"] = in.Status
	}

	if err := l.svc.DB.Model(&model.User{}).Where("id = ?", in.Id).Updates(updates).Error; err != nil {
		return &userv1.UserResponse{Code: 500, Message: err.Error()}, nil
	}

	return NewGetUserLogic(l.ctx, l.svc).GetUser(&userv1.GetUserRequest{Id: in.Id})
}
