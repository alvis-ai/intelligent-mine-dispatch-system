package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterVehicleRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := vehiclev1.NewVehicleServiceClient(conn)

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/vehicles",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			resp, _ := client.ListVehicle(r.Context(), &vehiclev1.ListVehicleRequest{Page: 1, PageSize: 100})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/vehicles",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req vehiclev1.CreateVehicleRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CreateVehicle(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/vehicles/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			var idUint uint64
			fmt.Sscanf(id, "%d", &idUint)
			resp, _ := client.GetVehicle(r.Context(), &vehiclev1.GetVehicleRequest{Id: idUint})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
