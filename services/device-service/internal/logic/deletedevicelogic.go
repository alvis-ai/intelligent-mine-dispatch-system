package logic

import (
	"context"

	devicev1 "github.com/aicong/mine-dispatch/proto/device/v1"
	"github.com/aicong/mine-dispatch/services/device-service/internal/model"
	"github.com/aicong/mine-dispatch/services/device-service/internal/svc"
)

type DeleteDeviceLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewDeleteDeviceLogic(ctx context.Context, svc *svc.ServiceContext) *DeleteDeviceLogic {
	return &DeleteDeviceLogic{ctx: ctx, svc: svc}
}

func (l *DeleteDeviceLogic) DeleteDevice(in *devicev1.DeleteDeviceRequest) (*devicev1.DeviceResponse, error) {
	if err := l.svc.DB.Delete(&model.Device{}, in.Id).Error; err != nil {
		return &devicev1.DeviceResponse{Code: 500, Message: err.Error()}, nil
	}
	return &devicev1.DeviceResponse{Code: 0, Message: "success"}, nil
}
