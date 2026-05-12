package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/model"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/svc"
)

type GetVehicleLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewGetVehicleLogic(ctx context.Context, svc *svc.ServiceContext) *GetVehicleLogic {
	return &GetVehicleLogic{ctx: ctx, svc: svc}
}

func (l *GetVehicleLogic) GetVehicle(in *vehiclev1.GetVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	var v model.Vehicle
	if err := l.svc.DB.First(&v, in.Id).Error; err != nil {
		return &vehiclev1.VehicleResponse{Code: 404, Message: "vehicle not found"}, nil
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
