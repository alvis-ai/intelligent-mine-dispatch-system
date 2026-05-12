package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/model"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/svc"
)

type DeleteVehicleLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewDeleteVehicleLogic(ctx context.Context, svc *svc.ServiceContext) *DeleteVehicleLogic {
	return &DeleteVehicleLogic{ctx: ctx, svc: svc}
}

func (l *DeleteVehicleLogic) DeleteVehicle(in *vehiclev1.DeleteVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	if err := l.svc.DB.Delete(&model.Vehicle{}, in.Id).Error; err != nil {
		return &vehiclev1.VehicleResponse{Code: 500, Message: err.Error()}, nil
	}
	return &vehiclev1.VehicleResponse{Code: 0, Message: "success"}, nil
}
