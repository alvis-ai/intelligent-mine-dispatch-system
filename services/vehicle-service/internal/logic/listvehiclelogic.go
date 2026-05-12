package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/model"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/svc"
)

type ListVehicleLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewListVehicleLogic(ctx context.Context, svc *svc.ServiceContext) *ListVehicleLogic {
	return &ListVehicleLogic{ctx: ctx, svc: svc}
}

func (l *ListVehicleLogic) ListVehicle(in *vehiclev1.ListVehicleRequest) (*vehiclev1.VehicleListResponse, error) {
	var vehicles []model.Vehicle
	var total int64
	db := l.svc.DB.Model(&model.Vehicle{})
	if in.Type != vehiclev1.VehicleType_VEHICLE_TYPE_UNSPECIFIED {
		db = db.Where("type = ?", in.Type)
	}
	if in.Status != vehiclev1.VehicleStatus_VEHICLE_STATUS_UNSPECIFIED {
		db = db.Where("status = ?", in.Status)
	}
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	db.Count(&total)
	if err := db.Offset(int((in.Page-1)*in.PageSize)).Limit(int(in.PageSize)).Find(&vehicles).Error; err != nil {
		return &vehiclev1.VehicleListResponse{Code: 500, Message: err.Error()}, nil
	}
	var list []*vehiclev1.Vehicle
	for _, v := range vehicles {
		list = append(list, &vehiclev1.Vehicle{
			Id: v.ID, Plate: v.Plate,
			Type: vehiclev1.VehicleType(v.Type),
			Status: vehiclev1.VehicleStatus(v.Status),
			Latitude: v.Latitude, Longitude: v.Longitude,
			FuelLevel: v.FuelLevel, MineId: v.MineID, DriverId: v.DriverID,
		})
	}
	return &vehiclev1.VehicleListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}
