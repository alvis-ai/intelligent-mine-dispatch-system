package server

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/telemetry/v1"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/svc"
)

type TelemetryServer struct {
	svc *svc.ServiceContext
	telemetryv1.UnimplementedTelemetryServiceServer
}

func NewTelemetryServer(svc *svc.ServiceContext) *TelemetryServer {
	return &TelemetryServer{svc: svc}
}

func (s *TelemetryServer) ReportLocation(ctx context.Context, in *telemetryv1.ReportLocationRequest) (*telemetryv1.ReportLocationResponse, error) {
	return logic.NewReportLocationLogic(ctx, s.svc).ReportLocation(in)
}

func (s *TelemetryServer) GetVehicleLocation(ctx context.Context, in *telemetryv1.GetVehicleLocationRequest) (*telemetryv1.GetVehicleLocationResponse, error) {
	return logic.NewGetVehicleLocationLogic(ctx, s.svc).GetVehicleLocation(in)
}

func (s *TelemetryServer) GetNearbyVehicles(ctx context.Context, in *telemetryv1.GetNearbyVehiclesRequest) (*telemetryv1.GetNearbyVehiclesResponse, error) {
	return logic.NewGetNearbyVehiclesLogic(ctx, s.svc).GetNearbyVehicles(in)
}
