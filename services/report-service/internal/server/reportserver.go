package server

import (
	"context"

	reportv1 "github.com/aicong/mine-dispatch/proto/report/v1"
	"github.com/aicong/mine-dispatch/services/report-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/report-service/internal/svc"
)

type ReportServer struct {
	svc *svc.ServiceContext
	reportv1.UnimplementedReportServiceServer
}

func NewReportServer(svc *svc.ServiceContext) *ReportServer {
	return &ReportServer{svc: svc}
}

func (s *ReportServer) GetDashboardSummary(ctx context.Context, in *reportv1.DashboardSummaryRequest) (*reportv1.DashboardSummaryResponse, error) {
	return logic.NewDashboardSummaryLogic(ctx, s.svc).GetDashboardSummary(in)
}

func (s *ReportServer) GetDispatchReport(ctx context.Context, in *reportv1.DispatchReportRequest) (*reportv1.DispatchReportResponse, error) {
	return logic.NewDispatchReportLogic(ctx, s.svc).GetDispatchReport(in)
}

func (s *ReportServer) GetVehicleUtilization(ctx context.Context, in *reportv1.VehicleUtilizationRequest) (*reportv1.VehicleUtilizationResponse, error) {
	return logic.NewVehicleUtilizationLogic(ctx, s.svc).GetVehicleUtilization(in)
}

func (s *ReportServer) GetTransportVolume(ctx context.Context, in *reportv1.TransportVolumeRequest) (*reportv1.TransportVolumeResponse, error) {
	return logic.NewTransportVolumeLogic(ctx, s.svc).GetTransportVolume(in)
}

func (s *ReportServer) GetAlarmTrend(ctx context.Context, in *reportv1.AlarmTrendRequest) (*reportv1.AlarmTrendResponse, error) {
	return logic.NewAlarmTrendLogic(ctx, s.svc).GetAlarmTrend(in)
}
