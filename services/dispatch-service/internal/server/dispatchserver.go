package server

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/dispatch/v1"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/svc"
)

type DispatchServer struct {
	svc *svc.ServiceContext
	dispatchv1.UnimplementedDispatchServiceServer
}

func NewDispatchServer(svc *svc.ServiceContext) *DispatchServer {
	return &DispatchServer{svc: svc}
}

func (s *DispatchServer) AssignTask(ctx context.Context, in *dispatchv1.AssignTaskRequest) (*dispatchv1.TaskResponse, error) {
	return logic.NewDispatchLogic(ctx, s.svc).AssignTask(in)
}

func (s *DispatchServer) GetTask(ctx context.Context, in *dispatchv1.GetTaskRequest) (*dispatchv1.TaskResponse, error) {
	return logic.NewDispatchLogic(ctx, s.svc).GetTask(in)
}

func (s *DispatchServer) CompleteTask(ctx context.Context, in *dispatchv1.CompleteTaskRequest) (*dispatchv1.TaskResponse, error) {
	return logic.NewDispatchLogic(ctx, s.svc).CompleteTask(in)
}

func (s *DispatchServer) CancelTask(ctx context.Context, in *dispatchv1.CancelTaskRequest) (*dispatchv1.TaskResponse, error) {
	return logic.NewDispatchLogic(ctx, s.svc).CancelTask(in)
}

func (s *DispatchServer) ListTask(ctx context.Context, in *dispatchv1.ListTaskRequest) (*dispatchv1.TaskListResponse, error) {
	return logic.NewDispatchLogic(ctx, s.svc).ListTask(in)
}
