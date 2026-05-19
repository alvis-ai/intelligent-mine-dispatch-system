package server

import (
	"context"

	aiv1 "github.com/aicong/mine-dispatch/proto/ai/v1"
	"github.com/aicong/mine-dispatch/services/ai-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/ai-service/internal/svc"
)

type AiServer struct {
	svc *svc.ServiceContext
	aiv1.UnimplementedAiServiceServer
}

func NewAiServer(svc *svc.ServiceContext) *AiServer {
	return &AiServer{svc: svc}
}

func (s *AiServer) PredictCongestion(ctx context.Context, in *aiv1.PredictCongestionRequest) (*aiv1.PredictCongestionResponse, error) {
	return logic.NewAiLogic(ctx, s.svc).PredictCongestion(in)
}

func (s *AiServer) RecommendRoute(ctx context.Context, in *aiv1.RecommendRouteRequest) (*aiv1.RecommendRouteResponse, error) {
	return logic.NewAiLogic(ctx, s.svc).RecommendRoute(in)
}

func (s *AiServer) PredictDemand(ctx context.Context, in *aiv1.PredictDemandRequest) (*aiv1.PredictDemandResponse, error) {
	return logic.NewAiLogic(ctx, s.svc).PredictDemand(in)
}

func (s *AiServer) SuggestAssign(ctx context.Context, in *aiv1.SuggestAssignRequest) (*aiv1.SuggestAssignResponse, error) {
	return logic.NewAiLogic(ctx, s.svc).SuggestAssign(in)
}

func (s *AiServer) HealthCheck(ctx context.Context, in *aiv1.HealthCheckRequest) (*aiv1.HealthCheckResponse, error) {
	return &aiv1.HealthCheckResponse{Status: "ok"}, nil
}
