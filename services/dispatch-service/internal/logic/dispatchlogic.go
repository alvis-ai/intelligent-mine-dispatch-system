package logic

import (
	"context"
	"fmt"

	"github.com/aicong/mine-dispatch/proto/dispatch/v1"
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
	task := model.DispatchTask{
		ID:          utils.NextID(),
		VehicleID:   vehicleID,
		LoadPointID: loadPointID,
		DumpPointID: dumpPointID,
		Status:      model.StatusActive,
		Algorithm:   "nearest_first",
	}
	if err := a.svc.DB.Create(&task).Error; err != nil {
		return 0, err
	}
	// Publish dispatch event
	msg := fmt.Sprintf(`{"task_id":%d,"vehicle_id":%d,"action":"assign"}`, task.ID, vehicleID)
	a.svc.Redis.Publish(a.ctx, "dispatch:events", msg)
	return task.ID, nil
}
