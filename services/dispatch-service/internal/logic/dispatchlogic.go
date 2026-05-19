package logic

import (
	"context"
	"fmt"
	"math"
	"time"

	dispatchv1 "github.com/aicong/mine-dispatch/proto/dispatch/v1"
	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/pkg/utils"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/model"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/svc"
)

type DispatchLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewDispatchLogic(ctx context.Context, svc *svc.ServiceContext) *DispatchLogic {
	return &DispatchLogic{ctx: ctx, svc: svc}
}

type TaskAssigner interface {
	Assign(vehicleID, loadPointID, dumpPointID uint64) (uint64, error)
}

func (l *DispatchLogic) getAssigner(algo string) TaskAssigner {
	switch algo {
	case "weighted_round_robin":
		return &WeightedRoundRobinAssigner{svc: l.svc}
	case "nearest_first":
		return &NearestFirstAssigner{svc: l.svc, ctx: l.ctx}
	case "genetic_algorithm":
		return NewGeneticAlgorithmAssigner(l.ctx, l.svc)
	case "ai_suggest":
		return NewAISuggestAssigner(l.ctx, l.svc)
	default:
		return &FIFOAssigner{svc: l.svc}
	}
}

func (l *DispatchLogic) AssignTask(in *dispatchv1.AssignTaskRequest) (*dispatchv1.TaskResponse, error) {
	assigner := l.getAssigner(in.Algorithm)
	taskID, err := assigner.Assign(in.VehicleId, in.LoadPointId, in.DumpPointId)
	if err != nil {
		return &dispatchv1.TaskResponse{Code: 500, Message: err.Error()}, nil
	}
	return l.GetTask(&dispatchv1.GetTaskRequest{Id: taskID})
}

func (l *DispatchLogic) GetTask(in *dispatchv1.GetTaskRequest) (*dispatchv1.TaskResponse, error) {
	var t model.DispatchTask
	if err := l.svc.DB.First(&t, in.Id).Error; err != nil {
		return &dispatchv1.TaskResponse{Code: 404, Message: "task not found"}, nil
	}
	return &dispatchv1.TaskResponse{
		Code: 0, Message: "success",
		Data: &dispatchv1.DispatchTask{
			Id:          t.ID,
			VehicleId:   t.VehicleID,
			LoadPointId: t.LoadPointID,
			DumpPointId: t.DumpPointID,
			Material:    t.Material,
			LoadLat:     t.LoadLat,
			LoadLon:     t.LoadLon,
			DumpLat:     t.DumpLat,
			DumpLon:     t.DumpLon,
			Status:      t.Status,
		},
	}, nil
}

func (l *DispatchLogic) CompleteTask(in *dispatchv1.CompleteTaskRequest) (*dispatchv1.TaskResponse, error) {
	l.svc.DB.Model(&model.DispatchTask{}).Where("id = ?", in.Id).Update("status", model.StatusCompleted)
	return l.GetTask(&dispatchv1.GetTaskRequest{Id: in.Id})
}

func (l *DispatchLogic) CancelTask(in *dispatchv1.CancelTaskRequest) (*dispatchv1.TaskResponse, error) {
	l.svc.DB.Model(&model.DispatchTask{}).Where("id = ?", in.Id).Update("status", model.StatusCancelled)
	return l.GetTask(&dispatchv1.GetTaskRequest{Id: in.Id})
}

func (l *DispatchLogic) ListTask(in *dispatchv1.ListTaskRequest) (*dispatchv1.TaskListResponse, error) {
	var tasks []model.DispatchTask
	var total int64
	db := l.svc.DB.Model(&model.DispatchTask{})
	if in.Status != "" {
		db = db.Where("status = ?", in.Status)
	}
	if in.VehicleId > 0 {
		db = db.Where("vehicle_id = ?", in.VehicleId)
	}
	db.Count(&total)
	if err := db.Order("created_at DESC").Offset(int((in.Page - 1) * in.PageSize)).Limit(int(in.PageSize)).Find(&tasks).Error; err != nil {
		return &dispatchv1.TaskListResponse{Code: 500, Message: err.Error()}, nil
	}
	var list []*dispatchv1.DispatchTask
	for _, t := range tasks {
		list = append(list, &dispatchv1.DispatchTask{
			Id: t.ID, VehicleId: t.VehicleID,
			LoadPointId: t.LoadPointID, DumpPointId: t.DumpPointID,
			Material: t.Material, Status: t.Status,
		})
	}
	return &dispatchv1.TaskListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}

// Assigners

type FIFOAssigner struct {
	svc *svc.ServiceContext
}

func (a *FIFOAssigner) Assign(vehicleID, loadPointID, dumpPointID uint64) (uint64, error) {
	task := model.DispatchTask{
		ID:          utils.NextID(),
		VehicleID:   vehicleID,
		LoadPointID: loadPointID,
		DumpPointID: dumpPointID,
		Status:      model.StatusActive,
		Algorithm:   "fifo",
	}
	if err := a.svc.DB.Create(&task).Error; err != nil {
		return 0, err
	}
	return task.ID, nil
}

type WeightedRoundRobinAssigner struct {
	svc *svc.ServiceContext
}

func (a *WeightedRoundRobinAssigner) Assign(vehicleID, loadPointID, dumpPointID uint64) (uint64, error) {
	task := model.DispatchTask{
		ID:          utils.NextID(),
		VehicleID:   vehicleID,
		LoadPointID: loadPointID,
		DumpPointID: dumpPointID,
		Status:      model.StatusActive,
		Algorithm:   "weighted_round_robin",
	}
	if err := a.svc.DB.Create(&task).Error; err != nil {
		return 0, err
	}
	return task.ID, nil
}

type NearestFirstAssigner struct {
	svc *svc.ServiceContext
	ctx context.Context
}

func (a *NearestFirstAssigner) Assign(vehicleID, loadPointID, dumpPointID uint64) (uint64, error) {
	// Query vehicle position
	var vhc struct {
		Latitude  float64
		Longitude float64
	}
	a.svc.DB.Table("vehicles").Select("latitude, longitude").Where("id = ?", vehicleID).Scan(&vhc)

	// Query loading point coordinates and material
	var loadPt LoadingPoint
	a.svc.DB.First(&loadPt, loadPointID)

	// Query dump point coordinates
	var dumpPt LoadingPoint
	a.svc.DB.First(&dumpPt, dumpPointID)

	// Calculate road distances via route-service
	loadDist, loadDur := a.getRouteDistance(vhc.Latitude, vhc.Longitude, loadPt.Latitude, loadPt.Longitude)
	dumpDist, dumpDur := a.getRouteDistance(loadPt.Latitude, loadPt.Longitude, dumpPt.Latitude, dumpPt.Longitude)
	totalDist := loadDist + dumpDist

	// Assign loading point material to task
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
		Algorithm:   "nearest_first",
	}
	if err := a.svc.DB.Create(&task).Error; err != nil {
		return 0, err
	}

	// Publish dispatch event with distance info
	msg := fmt.Sprintf(
		`{"task_id":%d,"vehicle_id":%d,"action":"assign","load_dist_m":%.0f,"dump_dist_m":%.0f,"total_dist_m":%.0f,"load_dur_s":%.0f,"dump_dur_s":%.0f}`,
		task.ID, vehicleID, loadDist, dumpDist, totalDist, loadDur, dumpDur,
	)
	a.svc.Redis.Publish(a.ctx, "dispatch:events", msg)
	return task.ID, nil
}

func (a *NearestFirstAssigner) getRouteDistance(fromLat, fromLon, toLat, toLon float64) (float64, float64) {
	if a.svc.RouteClient == nil {
		d := haversine(fromLat, fromLon, toLat, toLon)
		return d, (d / 1000) / 30 * 3600
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	resp, err := a.svc.RouteClient.GetDistance(ctx, &routev1.GetDistanceRequest{
		FromLat: fromLat,
		FromLon: fromLon,
		ToLat:   toLat,
		ToLon:   toLon,
	})
	if err != nil || resp == nil || resp.Code != 0 {
		d := haversine(fromLat, fromLon, toLat, toLon)
		return d, (d / 1000) / 30 * 3600
	}
	return resp.DistanceM, resp.DurationS
}

// haversine computes great-circle distance in meters between two lat/lon points.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
