package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aicong/mine-dispatch/proto/telemetry/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterTelemetryRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := telemetryv1.NewTelemetryServiceClient(conn)

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/telemetry/location",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req telemetryv1.ReportLocationRequest
			json.Unmarshal(body, &req)
			resp, _ := client.ReportLocation(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/telemetry/nearby",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			resp, _ := client.GetNearbyVehicles(r.Context(), &telemetryv1.GetNearbyVehiclesRequest{
				Latitude: 39.9, Longitude: 116.4, RadiusKm: 5,
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
