package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aicong/mine-dispatch/proto/telemetry/v1"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/svc"
)

type GetVehicleLocationLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewGetVehicleLocationLogic(ctx context.Context, svc *svc.ServiceContext) *GetVehicleLocationLogic {
	return &GetVehicleLocationLogic{ctx: ctx, svc: svc}
}

func (l *GetVehicleLocationLogic) GetVehicleLocation(in *telemetryv1.GetVehicleLocationRequest) (*telemetryv1.GetVehicleLocationResponse, error) {
	data, err := l.svc.Redis.Get(l.ctx, fmt.Sprintf("vehicle:loc:%d", in.VehicleId)).Bytes()
	if err != nil {
		return &telemetryv1.GetVehicleLocationResponse{Code: 404, Message: "location not found"}, nil
	}

	var locMap map[string]interface{}
	json.Unmarshal(data, &locMap)

	return &telemetryv1.GetVehicleLocationResponse{
		Code:    0,
		Message: "success",
		Location: &telemetryv1.LocationData{
			VehicleId: in.VehicleId,
			Latitude:  locMap["latitude"].(float64),
			Longitude: locMap["longitude"].(float64),
			Speed:     locMap["speed"].(float64),
			Heading:   locMap["heading"].(float64),
		},
	}, nil
}
