package logic

import (
	"context"
	"fmt"

	"github.com/aicong/mine-dispatch/proto/telemetry/v1"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/svc"
	"github.com/redis/go-redis/v9"
)

type GetNearbyVehiclesLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewGetNearbyVehiclesLogic(ctx context.Context, svc *svc.ServiceContext) *GetNearbyVehiclesLogic {
	return &GetNearbyVehiclesLogic{ctx: ctx, svc: svc}
}

func (l *GetNearbyVehiclesLogic) GetNearbyVehicles(in *telemetryv1.GetNearbyVehiclesRequest) (*telemetryv1.GetNearbyVehiclesResponse, error) {
	res, err := l.svc.Redis.GeoRadius(l.ctx, "vehicle:geo", in.Longitude, in.Latitude, &redis.GeoRadiusQuery{
		Radius:    in.RadiusKm,
		Unit:      "km",
		WithCoord: true,
		WithDist:  true,
	}).Result()
	if err != nil {
		return &telemetryv1.GetNearbyVehiclesResponse{Code: 500, Message: err.Error()}, nil
	}

	var vehicles []*telemetryv1.NearbyVehicle
	for _, v := range res {
		vehicles = append(vehicles, &telemetryv1.NearbyVehicle{
			VehicleId:  parseUint64(v.Name),
			Latitude:   v.Latitude,
			Longitude:  v.Longitude,
			DistanceKm: v.Dist,
		})
	}

	return &telemetryv1.GetNearbyVehiclesResponse{
		Code:     0,
		Message:  "success",
		Vehicles: vehicles,
	}, nil
}

func parseUint64(s string) uint64 {
	var id uint64
	fmt.Sscanf(s, "%d", &id)
	return id
}
