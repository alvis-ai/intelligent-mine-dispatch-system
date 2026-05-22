package handler

import (
	"encoding/json"
	"net/http"

	reportv1 "github.com/aicong/mine-dispatch/proto/report/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterReportRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := reportv1.NewReportServiceClient(conn)

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/reports/dashboard-summary",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			resp, _ := client.GetDashboardSummary(r.Context(), &reportv1.DashboardSummaryRequest{})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/reports/dispatch",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			resp, _ := client.GetDispatchReport(r.Context(), &reportv1.DispatchReportRequest{
				StartDate: q.Get("start_date"),
				EndDate:   q.Get("end_date"),
				GroupBy:   q.Get("group_by"),
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/reports/vehicle-utilization",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			resp, _ := client.GetVehicleUtilization(r.Context(), &reportv1.VehicleUtilizationRequest{
				StartDate: q.Get("start_date"),
				EndDate:   q.Get("end_date"),
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/reports/transport-volume",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			resp, _ := client.GetTransportVolume(r.Context(), &reportv1.TransportVolumeRequest{
				StartDate: q.Get("start_date"),
				EndDate:   q.Get("end_date"),
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/reports/alarm-trend",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			resp, _ := client.GetAlarmTrend(r.Context(), &reportv1.AlarmTrendRequest{
				StartDate: q.Get("start_date"),
				EndDate:   q.Get("end_date"),
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
