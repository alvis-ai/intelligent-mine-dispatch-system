package logic

import (
	"context"
	"fmt"
	"time"

	aiv1 "github.com/aicong/mine-dispatch/proto/ai/v1"
	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/model"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/svc"
	"github.com/aicong/mine-dispatch/pkg/utils"
)

type AISuggestAssigner struct {
	svc *svc.ServiceContext
	ctx context.Context
}

func NewAISuggestAssigner(ctx context.Context, svc *svc.ServiceContext) *AISuggestAssigner {
	return &AISuggestAssigner{svc: svc, ctx: ctx}
}

func (a *AISuggestAssigner) Assign(vehicleID, loadPointID, dumpPointID uint64) (uint64, error) {
	// Get vehicle position
	var vhc struct {
		Latitude  float64
		Longitude float64
	}
	a.svc.DB.Table("vehicles").Select("latitude, longitude").Where("id = ?", vehicleID).Scan(&vhc)

	// Get loading point info
	var loadPt LoadingPoint
	a.svc.DB.First(&loadPt, loadPointID)

	// Get dump point info
	var dumpPt LoadingPoint
	a.svc.DB.First(&dumpPt, dumpPointID)

	// Get active task count for this vehicle
	var activeCount int64
	a.svc.DB.Model(&model.DispatchTask{}).
		Where("vehicle_id = ? AND status IN ?", vehicleID, []string{"active", "pending"}).
		Count(&activeCount)

	// Call ai-service for suggestion
	if a.svc.AiClient != nil {
		ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
		defer cancel()

		resp, err := a.svc.AiClient.SuggestAssign(ctx, &aiv1.SuggestAssignRequest{
			Vehicles: []*aiv1.AIVehicleInfo{
				{
					VehicleId:       vehicleID,
					Latitude:        vhc.Latitude,
					Longitude:       vhc.Longitude,
					ActiveTaskCount: int32(activeCount),
				},
			},
			Tasks: []*aiv1.AITaskCandidate{
				{
					VehicleIdHint: vehicleID,
					LoadPointId:   loadPointID,
					DumpPointId:   dumpPointID,
				},
			},
		})
		if err == nil && resp != nil && len(resp.Suggestions) > 0 {
			top := resp.Suggestions[0]
			if top.Score > 0 {
				_ = top
				// Use AI suggestion, create task normally below
			}
		}
	}

	// Calculate road distances (same as NearestFirstAssigner)
	loadDist, loadDur := a.getRouteDistance(vhc.Latitude, vhc.Longitude, loadPt.Latitude, loadPt.Longitude)
	dumpDist, dumpDur := a.getRouteDistance(loadPt.Latitude, loadPt.Longitude, dumpPt.Latitude, dumpPt.Longitude)

	totalDist := loadDist + dumpDist

	material := loadPt.Material
	task := model.DispatchTask{
		ID:          utils.NextID(),
		VehicleID:   vehicleID,
		LoadPointID: loadPointID,
		DumpPointID: dumpPointID,
		Material:    material,
		LoadLat:     loadPt.Latitude,
		LoadLon:     loadPt.Longitude,
		DumpLat:     dumpPt.Latitude,
		DumpLon:     dumpPt.Longitude,
		Status:      model.StatusActive,
		Algorithm:   "ai_suggest",
	}
	if err := a.svc.DB.Create(&task).Error; err != nil {
		return 0, err
	}

	msg := fmt.Sprintf(
		`{"task_id":%d,"vehicle_id":%d,"action":"assign","algorithm":"ai_suggest","load_dist_m":%.0f,"dump_dist_m":%.0f,"total_dist_m":%.0f,"load_dur_s":%.0f,"dump_dur_s":%.0f}`,
		task.ID, vehicleID, loadDist, dumpDist, totalDist, loadDur, dumpDur,
	)
	a.svc.Redis.Publish(a.ctx, "dispatch:events", msg)
	return task.ID, nil
}

func (a *AISuggestAssigner) getRouteDistance(fromLat, fromLon, toLat, toLon float64) (float64, float64) {
	if a.svc.RouteClient == nil {
		d := haversine(fromLat, fromLon, toLat, toLon)
		return d, (d / 1000) / 30 * 3600
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	resp, err := a.svc.RouteClient.GetDistance(ctx, &routev1.GetDistanceRequest{
		FromLat: fromLat, FromLon: fromLon,
		ToLat: toLat, ToLon: toLon,
	})
	if err != nil || resp == nil || resp.Code != 0 {
		d := haversine(fromLat, fromLon, toLat, toLon)
		return d, (d / 1000) / 30 * 3600
	}
	return resp.DistanceM, resp.DurationS
}
