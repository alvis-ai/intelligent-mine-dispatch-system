package logic

import (
	"context"

	devicev1 "github.com/aicong/mine-dispatch/proto/device/v1"
	"github.com/aicong/mine-dispatch/services/device-service/internal/model"
	"github.com/aicong/mine-dispatch/services/device-service/internal/svc"
	"github.com/aicong/mine-dispatch/pkg/utils"
)

type CreateDeviceLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewCreateDeviceLogic(ctx context.Context, svc *svc.ServiceContext) *CreateDeviceLogic {
	return &CreateDeviceLogic{ctx: ctx, svc: svc}
}

func (l *CreateDeviceLogic) CreateDevice(in *devicev1.CreateDeviceRequest) (*devicev1.DeviceResponse, error) {
	d := model.Device{
		ID:         utils.NextID(),
		Name:       in.Name,
		DeviceType: int32(in.DeviceType),
		Status:     int32(devicev1.DeviceStatus_DEVICE_STATUS_ONLINE),
		MineID:     in.MineId,
		VehicleID:  in.VehicleId,
	}
	if err := l.svc.DB.Create(&d).Error; err != nil {
		return &devicev1.DeviceResponse{Code: 500, Message: err.Error()}, nil
	}
	return &devicev1.DeviceResponse{
		Code: 0, Message: "success",
		Data: deviceToProto(&d),
	}, nil
}
