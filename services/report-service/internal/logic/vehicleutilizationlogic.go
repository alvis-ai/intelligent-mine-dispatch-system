package logic

import (
	"context"

	reportv1 "github.com/aicong/mine-dispatch/proto/report/v1"
	"github.com/aicong/mine-dispatch/services/report-service/internal/svc"
)

type VehicleUtilizationLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewVehicleUtilizationLogic(ctx context.Context, svc *svc.ServiceContext) *VehicleUtilizationLogic {
	return &VehicleUtilizationLogic{ctx: ctx, svc: svc}
}

func (l *VehicleUtilizationLogic) GetVehicleUtilization(in *reportv1.VehicleUtilizationRequest) (*reportv1.VehicleUtilizationResponse, error) {
	var rows []*reportv1.VehicleUtilRow
	l.svc.DB.Raw(`
		SELECT v.id AS vehicle_id, v.plate,
			COALESCE(t.total_tasks, 0) AS total_tasks,
			COALESCE(t.completed_tasks, 0) AS completed_tasks,
			CASE WHEN COALESCE(t.total_tasks, 0) > 0
				THEN ROUND(COALESCE(t.completed_tasks, 0)::numeric / t.total_tasks::numeric, 4)
				ELSE 0
			END AS utilization_rate
		FROM vehicles v
		LEFT JOIN (
			SELECT vehicle_id,
				COUNT(*) AS total_tasks,
				COUNT(*) FILTER (WHERE status = 'completed') AS completed_tasks
			FROM dispatch_tasks
			WHERE (created_at::date >= ? AND created_at::date <= ?) OR (? = '' OR ? = '')
			GROUP BY vehicle_id
		) t ON v.id = t.vehicle_id
		WHERE v.mine_id = ? OR ? = 0
		ORDER BY total_tasks DESC
	`, in.StartDate, in.EndDate, in.StartDate, in.EndDate, in.MineId, in.MineId).Scan(&rows)
	if rows == nil {
		rows = []*reportv1.VehicleUtilRow{}
	}
	return &reportv1.VehicleUtilizationResponse{Rows: rows}, nil
}
