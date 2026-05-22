package logic

import (
	"context"

	reportv1 "github.com/aicong/mine-dispatch/proto/report/v1"
	"github.com/aicong/mine-dispatch/services/report-service/internal/svc"
)

type AlarmTrendLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewAlarmTrendLogic(ctx context.Context, svc *svc.ServiceContext) *AlarmTrendLogic {
	return &AlarmTrendLogic{ctx: ctx, svc: svc}
}

func (l *AlarmTrendLogic) GetAlarmTrend(in *reportv1.AlarmTrendRequest) (*reportv1.AlarmTrendResponse, error) {
	var rows []*reportv1.AlarmTrendRow
	l.svc.DB.Raw(`
		SELECT created_at::date::text AS date,
			severity,
			COUNT(*) AS count
		FROM alarm_events
		WHERE (created_at::date >= ? AND created_at::date <= ?) OR (? = '' OR ? = '')
		GROUP BY created_at::date, severity
		ORDER BY date, severity
	`, in.StartDate, in.EndDate, in.StartDate, in.EndDate).Scan(&rows)
	if rows == nil {
		rows = []*reportv1.AlarmTrendRow{}
	}
	return &reportv1.AlarmTrendResponse{Rows: rows}, nil
}
