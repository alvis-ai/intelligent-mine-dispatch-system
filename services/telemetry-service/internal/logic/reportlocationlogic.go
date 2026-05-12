package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aicong/mine-dispatch/proto/telemetry/v1"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/svc"
	"github.com/redis/go-redis/v9"
)

type ReportLocationLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewReportLocationLogic(ctx context.Context, svc *svc.ServiceContext) *ReportLocationLogic {
	return &ReportLocationLogic{ctx: ctx, svc: svc}
}

const (
	geoKey       = "vehicle:geo"
	locationKey  = "vehicle:loc:%d"
	locationChan = "vehicle:location:updates"
)

func (l *ReportLocationLogic) ReportLocation(in *telemetryv1.ReportLocationRequest) (*telemetryv1.ReportLocationResponse, error) {
	loc := in.Location
	vid := loc.VehicleId

	// Store in Redis GEO set
	l.svc.Redis.GeoAdd(l.ctx, geoKey, &redis.GeoLocation{
		Name:      fmt.Sprintf("%d", vid),
		Longitude: loc.Longitude,
		Latitude:  loc.Latitude,
	})

	// Store latest location with TTL
	locData, _ := json.Marshal(map[string]interface{}{
		"vehicle_id": vid,
		"latitude":   loc.Latitude,
		"longitude":  loc.Longitude,
		"altitude":   loc.Altitude,
		"speed":      loc.Speed,
		"heading":    loc.Heading,
		"timestamp":  time.Now().UnixMilli(),
	})
	l.svc.Redis.Set(l.ctx, fmt.Sprintf(locationKey, vid), locData, 30*time.Second)

	// Publish for WebSocket
	l.svc.Redis.Publish(l.ctx, locationChan, string(locData))

	return &telemetryv1.ReportLocationResponse{Code: 0, Message: "success"}, nil
}
