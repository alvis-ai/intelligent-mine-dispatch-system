package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/pkg/utils"
	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/model"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/svc"
)

type CreateVehicleLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewCreateVehicleLogic(ctx context.Context, svc *svc.ServiceContext) *CreateVehicleLogic {
	return &CreateVehicleLogic{ctx: ctx, svc: svc}
}

func (l *CreateVehicleLogic) CreateVehicle(in *vehiclev1.CreateVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	v := model.Vehicle{
		ID:     utils.NextID(),
		Plate:  in.Plate,
		Type:   int32(in.Type),
		Status: int32(vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE),
		MineID: in.MineId,
	}
	if err := l.svc.DB.Create(&v).Error; err != nil {
		return &vehiclev1.VehicleResponse{Code: 500, Message: err.Error()}, nil
	}
	return &vehiclev1.VehicleResponse{
		Code: 0, Message: "success",
		Data: &vehiclev1.Vehicle{
			Id: v.ID, Plate: v.Plate,
			Type: vehiclev1.VehicleType(v.Type),
			Status: vehiclev1.VehicleStatus(v.Status),
			Latitude: v.Latitude, Longitude: v.Longitude,
			FuelLevel: v.FuelLevel, MineId: v.MineID, DriverId: v.DriverID,
		},
	}, nil
}
