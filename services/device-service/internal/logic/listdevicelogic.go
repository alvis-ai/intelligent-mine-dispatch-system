package logic

import (
	"context"
	"time"

	devicev1 "github.com/aicong/mine-dispatch/proto/device/v1"
	"github.com/aicong/mine-dispatch/services/device-service/internal/model"
	"github.com/aicong/mine-dispatch/services/device-service/internal/svc"
)

type ListDeviceLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewListDeviceLogic(ctx context.Context, svc *svc.ServiceContext) *ListDeviceLogic {
	return &ListDeviceLogic{ctx: ctx, svc: svc}
}

func (l *ListDeviceLogic) ListDevice(in *devicev1.ListDeviceRequest) (*devicev1.DeviceListResponse, error) {
	var devices []model.Device
	var total int64
	db := l.svc.DB.Model(&model.Device{})
	if in.DeviceType != devicev1.DeviceType_DEVICE_TYPE_UNSPECIFIED {
		db = db.Where("device_type = ?", in.DeviceType)
	}
	if in.Status != devicev1.DeviceStatus_DEVICE_STATUS_UNSPECIFIED {
		db = db.Where("status = ?", in.Status)
	}
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	if in.VehicleId > 0 {
		db = db.Where("vehicle_id = ?", in.VehicleId)
	}
	db.Count(&total)
	page := int(in.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize < 1 {
		pageSize = 20
	}
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&devices).Error; err != nil {
		return &devicev1.DeviceListResponse{Code: 500, Message: err.Error()}, nil
	}
	var list []*devicev1.Device
	for i := range devices {
		list = append(list, deviceToProto(&devices[i]))
	}
	return &devicev1.DeviceListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}

func deviceToProto(d *model.Device) *devicev1.Device {
	lastOnline := ""
	if !d.LastOnlineAt.IsZero() {
		lastOnline = d.LastOnlineAt.Format(time.RFC3339)
	}
	return &devicev1.Device{
		Id:              d.ID,
		Name:            d.Name,
		DeviceType:      devicev1.DeviceType(d.DeviceType),
		Status:          devicev1.DeviceStatus(d.Status),
		FirmwareVersion: d.FirmwareVersion,
		Latitude:        d.Latitude,
		Longitude:       d.Longitude,
		MineId:          d.MineID,
		VehicleId:       d.VehicleID,
		LastOnlineAt:    lastOnline,
		CreatedAt:       d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       d.UpdatedAt.Format(time.RFC3339),
	}
}
