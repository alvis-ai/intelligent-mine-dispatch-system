package logic

import (
	"context"

	reportv1 "github.com/aicong/mine-dispatch/proto/report/v1"
	"github.com/aicong/mine-dispatch/services/report-service/internal/svc"
)

type DispatchReportLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewDispatchReportLogic(ctx context.Context, svc *svc.ServiceContext) *DispatchReportLogic {
	return &DispatchReportLogic{ctx: ctx, svc: svc}
}

func (l *DispatchReportLogic) GetDispatchReport(in *reportv1.DispatchReportRequest) (*reportv1.DispatchReportResponse, error) {
	var rows []*reportv1.DispatchReportRow

	query := `SELECT %s AS dimension,
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE status = 'completed') AS completed,
		COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
		COUNT(*) FILTER (WHERE status = 'active') AS active,
		COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at)) / 60), 0) AS avg_duration_minutes
	FROM dispatch_tasks
	WHERE (created_at::date >= ? AND created_at::date <= ?) OR (? = '' OR ? = '')
	GROUP BY dimension ORDER BY dimension`

	switch in.GroupBy {
	case "algorithm":
		query = `SELECT algorithm AS dimension,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'completed') AS completed,
			COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
			COUNT(*) FILTER (WHERE status = 'active') AS active,
			COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at)) / 60), 0) AS avg_duration_minutes
		FROM dispatch_tasks
		WHERE (created_at::date >= ? AND created_at::date <= ?) OR (? = '' OR ? = '')
		GROUP BY algorithm ORDER BY dimension`
	case "status":
		query = `SELECT status AS dimension,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'completed') AS completed,
			COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
			COUNT(*) FILTER (WHERE status = 'active') AS active,
			COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at)) / 60), 0) AS avg_duration_minutes
		FROM dispatch_tasks
		WHERE (created_at::date >= ? AND created_at::date <= ?) OR (? = '' OR ? = '')
		GROUP BY status ORDER BY dimension`
	default:
		query = `SELECT created_at::date::text AS dimension,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'completed') AS completed,
			COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
			COUNT(*) FILTER (WHERE status = 'active') AS active,
			COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at)) / 60), 0) AS avg_duration_minutes
		FROM dispatch_tasks
		WHERE (created_at::date >= ? AND created_at::date <= ?) OR (? = '' OR ? = '')
		GROUP BY created_at::date ORDER BY dimension`
	}

	l.svc.DB.Raw(query, in.StartDate, in.EndDate, in.StartDate, in.EndDate).Scan(&rows)
	if rows == nil {
		rows = []*reportv1.DispatchReportRow{}
	}
	return &reportv1.DispatchReportResponse{Rows: rows}, nil
}
