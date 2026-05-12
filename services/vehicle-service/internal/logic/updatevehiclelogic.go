package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/model"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/svc"
)

type UpdateVehicleLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewUpdateVehicleLogic(ctx context.Context, svc *svc.ServiceContext) *UpdateVehicleLogic {
	return &UpdateVehicleLogic{ctx: ctx, svc: svc}
}

func (l *UpdateVehicleLogic) UpdateVehicle(in *vehiclev1.UpdateVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	updates := map[string]interface{}{}
	if in.Status != vehiclev1.VehicleStatus_VEHICLE_STATUS_UNSPECIFIED {
		updates["status"] = int32(in.Status)
	}
	if in.DriverId > 0 {
		updates["driver_id"] = in.DriverId
	}
	if len(updates) > 0 {
		l.svc.DB.Model(&model.Vehicle{}).Where("id = ?", in.Id).Updates(updates)
	}
	return NewGetVehicleLogic(l.ctx, l.svc).GetVehicle(&vehiclev1.GetVehicleRequest{Id: in.Id})
}
