package logic

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	aiv1 "github.com/aicong/mine-dispatch/proto/ai/v1"
	"github.com/aicong/mine-dispatch/services/ai-service/internal/model"
)

func (l *AiLogic) PredictCongestion(in *aiv1.PredictCongestionRequest) (*aiv1.PredictCongestionResponse, error) {
	mineID := in.MineId
	if mineID == 0 {
		mineID = 1
	}
	lookback := in.LookbackMinutes
	if lookback <= 0 {
		lookback = 60
	}

	// Try Redis cache
	cacheKey := fmt.Sprintf("ai:congestion:%d:%d", mineID, lookback)
	cached, err := l.svc.Redis.Get(l.ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var resp aiv1.PredictCongestionResponse
		if json.Unmarshal([]byte(cached), &resp) == nil {
			return &resp, nil
		}
	}

	// Load data
	var edges []model.RoadEdge
	var nodes []model.RoadNode
	var pts []model.LoadingPoint
	var tasks []model.DispatchTask

	l.svc.DB.Where("mine_id = ?", mineID).Find(&edges)
	l.svc.DB.Where("mine_id = ?", mineID).Find(&nodes)
	l.svc.DB.Where("mine_id = ? OR mine_id = 0", mineID).Find(&pts)

	since := time.Now().Add(-time.Duration(lookback) * time.Minute)
	l.svc.DB.Where("updated_at > ? AND status IN ?", since, []string{"active", "completed"}).
		Find(&tasks)

	if len(nodes) == 0 {
		return &aiv1.PredictCongestionResponse{Code: 0, Message: "no road data", Data: []*aiv1.EdgeCongestion{}}, nil
	}

	// Build node position map
	nodePos := make(map[uint64]*model.RoadNode)
	for i := range nodes {
		nodePos[nodes[i].ID] = &nodes[i]
	}

	// Map loading_points to nearest road node
	lpNode := make(map[uint64]uint64) // loading_point_id → nearest node_id
	for _, pt := range pts {
		nearest := findNearestNodeID(pt.Latitude, pt.Longitude, nodes)
		if nearest > 0 {
			lpNode[pt.ID] = nearest
		}
	}

	// Count traffic per node from recent tasks
	nodeTraffic := make(map[uint64]int)
	for _, t := range tasks {
		if nid, ok := lpNode[t.LoadPointID]; ok {
			nodeTraffic[nid]++
		}
		if nid, ok := lpNode[t.DumpPointID]; ok {
			nodeTraffic[nid]++
		}
	}

	// Compute congestion per edge
	totalTasks := len(tasks)
	var results []*aiv1.EdgeCongestion
	for _, e := range edges {
		traffic := nodeTraffic[e.FromNodeID] + nodeTraffic[e.ToNodeID]
		capacity := math.Max(1.0, float64(e.MaxSpeedKMH)*0.15)
		score := math.Min(1.0, float64(traffic)/capacity)

		confidence := math.Min(1.0, float64(totalTasks)/20.0)
		predictedSpeed := float64(e.MaxSpeedKMH) * (1 - score*0.4)

		results = append(results, &aiv1.EdgeCongestion{
			EdgeId:              e.ID,
			FromNodeId:          e.FromNodeID,
			ToNodeId:            e.ToNodeID,
			CongestionScore:     math.Round(score*100) / 100,
			PredictedSpeedKmh:   math.Round(predictedSpeed*10) / 10,
			PredictedVehicleCount: int32(traffic),
			Confidence:           math.Round(confidence*100) / 100,
		})
	}

	resp := &aiv1.PredictCongestionResponse{
		Code: 0,
		Data: results,
	}

	// Cache 60s in Redis
	if data, err := json.Marshal(resp); err == nil {
		l.svc.Redis.Set(l.ctx, cacheKey, string(data), 60*time.Second)
	}

	return resp, nil
}

func findNearestNodeID(lat, lon float64, nodes []model.RoadNode) uint64 {
	var bestID uint64
	bestDist := math.MaxFloat64
	for _, n := range nodes {
		d := haversineAI(lat, lon, n.Latitude, n.Longitude)
		if d < bestDist {
			bestDist = d
			bestID = n.ID
		}
	}
	return bestID
}

func haversineAI(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
