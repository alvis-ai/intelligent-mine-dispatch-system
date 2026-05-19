package logic

import (
	"math"
	"time"

	aiv1 "github.com/aicong/mine-dispatch/proto/ai/v1"
	"github.com/aicong/mine-dispatch/services/ai-service/internal/model"
)

func (l *AiLogic) PredictDemand(in *aiv1.PredictDemandRequest) (*aiv1.PredictDemandResponse, error) {
	mineID := in.MineId
	if mineID == 0 {
		mineID = 1
	}

	// Load all loading/dumping points
	var pts []model.LoadingPoint
	l.svc.DB.Where("(mine_id = ? OR mine_id = 0) AND status = 1", mineID).Find(&pts)

	if len(pts) == 0 {
		return &aiv1.PredictDemandResponse{Code: 0, Message: "no points found", Data: []*aiv1.LoadingPointDemand{}}, nil
	}

	now := time.Now()
	hourAgo := now.Add(-1 * time.Hour)
	sixHoursAgo := now.Add(-6 * time.Hour)
	dayAgo := now.Add(-24 * time.Hour)

	var results []*aiv1.LoadingPointDemand
	for _, pt := range pts {
		// Count tasks per time window
		var recent1h, recent6h, recent24h int64
		var pendingCount, activeCount int64

		l.svc.DB.Model(&model.DispatchTask{}).
			Where("load_point_id = ? AND created_at > ?", pt.ID, hourAgo).Count(&recent1h)
		l.svc.DB.Model(&model.DispatchTask{}).
			Where("load_point_id = ? AND created_at > ?", pt.ID, sixHoursAgo).Count(&recent6h)
		l.svc.DB.Model(&model.DispatchTask{}).
			Where("load_point_id = ? AND created_at > ?", pt.ID, dayAgo).Count(&recent24h)
		l.svc.DB.Model(&model.DispatchTask{}).
			Where("load_point_id = ? AND status = 'pending'", pt.ID).Count(&pendingCount)
		l.svc.DB.Model(&model.DispatchTask{}).
			Where("load_point_id = ? AND status = 'active'", pt.ID).Count(&activeCount)

		// EWMA-like scoring: recent traffic has higher weight
		trafficScore := float64(recent1h)*0.6 + float64(recent6h-recent1h)*0.3 + float64(recent24h-recent6h)*0.1
		normalizedTraffic := math.Min(1.0, trafficScore/20.0)

		pendingScore := math.Min(1.0, float64(pendingCount)/10.0)
		activeScore := math.Min(1.0, float64(activeCount)/10.0)

		demandScore := normalizedTraffic*0.5 + pendingScore*0.3 + activeScore*0.2
		demandScore = math.Round(demandScore*100) / 100

		dataPoints := recent24h + pendingCount + activeCount
		confidence := math.Min(1.0, float64(dataPoints)/30.0)
		confidence = math.Round(confidence*100) / 100

		results = append(results, &aiv1.LoadingPointDemand{
			LoadPointId:       pt.ID,
			Name:              pt.Name,
			PointType:         pt.Type,
			Material:          pt.Material,
			DemandScore:       demandScore,
			PendingTaskCount:  int32(pendingCount),
			ActiveVehicleCount: int32(activeCount),
			Confidence:        confidence,
		})
	}

	return &aiv1.PredictDemandResponse{Code: 0, Message: "success", Data: results}, nil
}
