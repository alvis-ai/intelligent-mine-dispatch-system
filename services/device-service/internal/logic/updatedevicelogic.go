package logic

import (
	"context"

	devicev1 "github.com/aicong/mine-dispatch/proto/device/v1"
	"github.com/aicong/mine-dispatch/services/device-service/internal/model"
	"github.com/aicong/mine-dispatch/services/device-service/internal/svc"
)

type UpdateDeviceLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewUpdateDeviceLogic(ctx context.Context, svc *svc.ServiceContext) *UpdateDeviceLogic {
	return &UpdateDeviceLogic{ctx: ctx, svc: svc}
}

func (l *UpdateDeviceLogic) UpdateDevice(in *devicev1.UpdateDeviceRequest) (*devicev1.DeviceResponse, error) {
	updates := map[string]interface{}{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Status != devicev1.DeviceStatus_DEVICE_STATUS_UNSPECIFIED {
		updates["status"] = int32(in.Status)
	}
	if in.FirmwareVersion != "" {
		updates["firmware_version"] = in.FirmwareVersion
	}
	if in.VehicleId > 0 {
		updates["vehicle_id"] = in.VehicleId
	}
	if len(updates) > 0 {
		l.svc.DB.Model(&model.Device{}).Where("id = ?", in.Id).Updates(updates)
	}
	return NewGetDeviceLogic(l.ctx, l.svc).GetDevice(&devicev1.GetDeviceRequest{Id: in.Id})
}
