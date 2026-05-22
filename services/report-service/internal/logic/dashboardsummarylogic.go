package logic

import (
	"context"
	"time"

	reportv1 "github.com/aicong/mine-dispatch/proto/report/v1"
	"github.com/aicong/mine-dispatch/services/report-service/internal/svc"
)

type DashboardSummaryLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewDashboardSummaryLogic(ctx context.Context, svc *svc.ServiceContext) *DashboardSummaryLogic {
	return &DashboardSummaryLogic{ctx: ctx, svc: svc}
}

func (l *DashboardSummaryLogic) GetDashboardSummary(in *reportv1.DashboardSummaryRequest) (*reportv1.DashboardSummaryResponse, error) {
	var totalVehicles int64
	l.svc.DB.Raw("SELECT COUNT(*) FROM vehicles WHERE mine_id = ? OR ? = 0", in.MineId, in.MineId).Scan(&totalVehicles)

	var activeTasks, pendingTasks, completedTasks int64
	l.svc.DB.Raw("SELECT COUNT(*) FROM dispatch_tasks WHERE status = 'active'").Scan(&activeTasks)
	l.svc.DB.Raw("SELECT COUNT(*) FROM dispatch_tasks WHERE status = 'pending'").Scan(&pendingTasks)
	l.svc.DB.Raw("SELECT COUNT(*) FROM dispatch_tasks WHERE status = 'completed'").Scan(&completedTasks)

	var unackCritical, unackWarning int64
	l.svc.DB.Raw("SELECT COUNT(*) FROM alarm_events WHERE acknowledged = false AND severity = 'critical'").Scan(&unackCritical)
	l.svc.DB.Raw("SELECT COUNT(*) FROM alarm_events WHERE acknowledged = false AND severity = 'warning'").Scan(&unackWarning)

	var todayDispatched int64
	today := time.Now().Format("2006-01-02")
	l.svc.DB.Raw("SELECT COUNT(*) FROM dispatch_tasks WHERE created_at::date = ?", today).Scan(&todayDispatched)

	return &reportv1.DashboardSummaryResponse{
		TotalVehicles:        int32(totalVehicles),
		ActiveTasks:          int32(activeTasks),
		PendingTasks:         int32(pendingTasks),
		CompletedTasks:       int32(completedTasks),
		UnacknowledgedCritical: int32(unackCritical),
		UnacknowledgedWarning:  int32(unackWarning),
		TodayDispatched:      int32(todayDispatched),
		LastUpdated:          time.Now().Format(time.RFC3339),
	}, nil
}
