package logic

import (
	"context"

	devicev1 "github.com/aicong/mine-dispatch/proto/device/v1"
	"github.com/aicong/mine-dispatch/services/device-service/internal/model"
	"github.com/aicong/mine-dispatch/services/device-service/internal/svc"
)

type GetDeviceLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewGetDeviceLogic(ctx context.Context, svc *svc.ServiceContext) *GetDeviceLogic {
	return &GetDeviceLogic{ctx: ctx, svc: svc}
}

func (l *GetDeviceLogic) GetDevice(in *devicev1.GetDeviceRequest) (*devicev1.DeviceResponse, error) {
	var d model.Device
	if err := l.svc.DB.First(&d, in.Id).Error; err != nil {
		return &devicev1.DeviceResponse{Code: 404, Message: "device not found"}, nil
	}
	return &devicev1.DeviceResponse{
		Code: 0, Message: "success",
		Data: deviceToProto(&d),
	}, nil
}
