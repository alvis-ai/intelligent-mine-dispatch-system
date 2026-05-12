package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/aicong/mine-dispatch/services/user-service/internal/model"
	"github.com/aicong/mine-dispatch/services/user-service/internal/svc"
)

type ListUserLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewListUserLogic(ctx context.Context, svc *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{ctx: ctx, svc: svc}
}

func (l *ListUserLogic) ListUser(in *userv1.ListUserRequest) (*userv1.UserListResponse, error) {
	var users []model.User
	var total int64

	db := l.svc.DB.Model(&model.User{})
	if in.Keyword != "" {
		db = db.Where("username LIKE ? OR real_name LIKE ?", "%"+in.Keyword+"%", "%"+in.Keyword+"%")
	}
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	db.Count(&total)

	offset := int((in.Page - 1) * in.PageSize)
	if err := db.Offset(offset).Limit(int(in.PageSize)).Find(&users).Error; err != nil {
		return &userv1.UserListResponse{Code: 500, Message: err.Error()}, nil
	}

	var list []*userv1.User
	for _, u := range users {
		list = append(list, &userv1.User{
			Id:       u.ID,
			Username: u.Username,
			Password: u.Password,
			RealName: u.RealName,
			Email:    u.Email,
			Phone:    u.Phone,
			Role:     u.Role,
			Status:   u.Status,
			MineId:   u.MineID,
		})
	}

	return &userv1.UserListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}
