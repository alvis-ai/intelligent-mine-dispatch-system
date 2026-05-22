package logic

import (
	"context"

	reportv1 "github.com/aicong/mine-dispatch/proto/report/v1"
	"github.com/aicong/mine-dispatch/services/report-service/internal/svc"
)

type TransportVolumeLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewTransportVolumeLogic(ctx context.Context, svc *svc.ServiceContext) *TransportVolumeLogic {
	return &TransportVolumeLogic{ctx: ctx, svc: svc}
}

func (l *TransportVolumeLogic) GetTransportVolume(in *reportv1.TransportVolumeRequest) (*reportv1.TransportVolumeResponse, error) {
	var rows []*reportv1.TransportVolumeRow
	l.svc.DB.Raw(`
		SELECT t.material,
			COUNT(*) AS task_count,
			t.load_point_id,
			COALESCE(lp.name, '未知') AS loading_point_name
		FROM dispatch_tasks t
		LEFT JOIN loading_points lp ON t.load_point_id = lp.id
		WHERE (t.created_at::date >= ? AND t.created_at::date <= ?) OR (? = '' OR ? = '')
		GROUP BY t.material, t.load_point_id, lp.name
		ORDER BY task_count DESC
	`, in.StartDate, in.EndDate, in.StartDate, in.EndDate).Scan(&rows)
	if rows == nil {
		rows = []*reportv1.TransportVolumeRow{}
	}
	return &reportv1.TransportVolumeResponse{Rows: rows}, nil
}
