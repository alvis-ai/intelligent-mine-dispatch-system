package logic

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	aiv1 "github.com/aicong/mine-dispatch/proto/ai/v1"
	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/services/ai-service/internal/model"
)

type matInfo struct {
	count int
	total int
}

func (l *AiLogic) SuggestAssign(in *aiv1.SuggestAssignRequest) (*aiv1.SuggestAssignResponse, error) {
	mineID := in.MineId
	if mineID == 0 {
		mineID = 1
	}

	var pts []model.LoadingPoint
	l.svc.DB.Where("mine_id = ? OR mine_id = 0", mineID).Find(&pts)
	ptMap := make(map[uint64]model.LoadingPoint)
	for _, pt := range pts {
		ptMap[pt.ID] = pt
	}

	congResp, _ := l.PredictCongestion(&aiv1.PredictCongestionRequest{
		MineId: mineID, LookbackMinutes: 30,
	})
	congMap := make(map[uint64]float64)
	if congResp != nil {
		for _, ec := range congResp.Data {
			congMap[ec.EdgeId] = ec.CongestionScore
		}
	}

	var activeTasks []struct {
		VehicleID uint64
		Count     int
	}
	l.svc.DB.Table("dispatch_tasks").
		Select("vehicle_id, COUNT(*) as count").
		Where("status = 'active' OR status = 'pending'").
		Group("vehicle_id").
		Scan(&activeTasks)
	activeMap := make(map[uint64]int)
	for _, at := range activeTasks {
		activeMap[at.VehicleID] = at.Count
	}

	var materialHistory []struct {
		VehicleID uint64
		Material  string
		TaskCount int
	}
	l.svc.DB.Table("dispatch_tasks").
		Select("vehicle_id, material, COUNT(*) as task_count").
		Where("material != ''").
		Group("vehicle_id, material").
		Scan(&materialHistory)
	matHistory := make(map[uint64]map[string]*matInfo)
	for _, mh := range materialHistory {
		if matHistory[mh.VehicleID] == nil {
			matHistory[mh.VehicleID] = make(map[string]*matInfo)
		}
		matHistory[mh.VehicleID][mh.Material] = &matInfo{count: mh.TaskCount}
	}
	for _, mats := range matHistory {
		total := 0
		for _, mi := range mats {
			total += mi.count
		}
		for _, mi := range mats {
			mi.total = total
		}
	}

	var suggestions []*aiv1.AISuggestionItem
	for _, task := range in.Tasks {
		loadPt, _ := ptMap[task.LoadPointId]
		dumpPt, _ := ptMap[task.DumpPointId]

		for _, veh := range in.Vehicles {
			loadDist, loadDur := l.getDist(veh.Latitude, veh.Longitude, task.LoadPointId, loadPt)
			dumpDist, dumpDur := l.getDist(loadPt.Latitude, loadPt.Longitude, task.DumpPointId, dumpPt)

			totalDist := loadDist + dumpDist
			totalDur := loadDur + dumpDur

			distScore := calcDistScore(totalDist)
			loadScore := calcLoadScore(activeMap[veh.VehicleId])
			congScore := calcCongScore(congMap)
			matScore := calcMatScore(matHistory[veh.VehicleId], loadPt.Material)
			utilScore := calcUtilScore(activeMap[veh.VehicleId], len(in.Vehicles))

			totalScore := distScore*0.40 + loadScore*0.25 + congScore*0.15 + matScore*0.10 + utilScore*0.10
			totalScore = math.Round(totalScore*1000) / 1000

			reason := fmt.Sprintf(
				"dist=%.2f load=%.2f cong=%.2f mat=%.2f util=%.2f",
				distScore, loadScore, congScore, matScore, utilScore,
			)

			suggestions = append(suggestions, &aiv1.AISuggestionItem{
				VehicleId:          veh.VehicleId,
				LoadPointId:        task.LoadPointId,
				DumpPointId:        task.DumpPointId,
				Score:              totalScore,
				EstimatedDistanceM: math.Round(totalDist*10) / 10,
				EstimatedDurationS: math.Round(totalDur*10) / 10,
				Reason:             reason,
			})
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	maxSuggestions := len(in.Vehicles) * len(in.Tasks)
	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}

	return &aiv1.SuggestAssignResponse{
		Code:        0,
		Message:     "success",
		Suggestions: suggestions,
	}, nil
}

func (l *AiLogic) getDist(fromLat, fromLon float64, pointID uint64, pt model.LoadingPoint) (float64, float64) {
	var toLat, toLon float64
	if pointID > 0 && pt.ID > 0 {
		toLat = pt.Latitude
		toLon = pt.Longitude
	} else {
		return 0, 0
	}

	if l.svc.RouteClient != nil {
		ctx, cancel := context.WithTimeout(l.ctx, 3*time.Second)
		defer cancel()
		resp, err := l.svc.RouteClient.GetDistance(ctx, &routev1.GetDistanceRequest{
			FromLat: fromLat, FromLon: fromLon,
			ToLat: toLat, ToLon: toLon,
		})
		if err == nil && resp != nil && resp.Code == 0 {
			return resp.DistanceM, resp.DurationS
		}
	}

	d := haversineAI(fromLat, fromLon, toLat, toLon)
	return d, (d / 1000) / 30 * 3600
}

func calcDistScore(distM float64) float64 {
	if distM <= 0 {
		return 1.0
	}
	return math.Max(0, 1.0-distM/20000.0)
}

func calcLoadScore(activeCount int) float64 {
	return math.Max(0, 1.0-float64(activeCount)*0.25)
}

func calcCongScore(congMap map[uint64]float64) float64 {
	if len(congMap) == 0 {
		return 0.5
	}
	total := 0.0
	for _, v := range congMap {
		total += v
	}
	return 1.0 - total/float64(len(congMap))
}

func calcMatScore(history map[string]*matInfo, material string) float64 {
	if history == nil || material == "" {
		return 0.5
	}
	if mi, ok := history[material]; ok && mi.total > 0 {
		return 0.5 + 0.5*float64(mi.count)/float64(mi.total)
	}
	return 0.3
}

func calcUtilScore(activeCount int, totalVehicles int) float64 {
	if totalVehicles <= 1 {
		return 0.5
	}
	return math.Max(0, 1.0-float64(activeCount)/float64(totalVehicles))
}
