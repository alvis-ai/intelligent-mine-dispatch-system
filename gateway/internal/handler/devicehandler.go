package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	devicev1 "github.com/aicong/mine-dispatch/proto/device/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterDeviceRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := devicev1.NewDeviceServiceClient(conn)

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/devices",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))
			pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
			if page < 1 {
				page = 1
			}
			if pageSize < 1 {
				pageSize = 20
			}
			mineID, _ := strconv.ParseUint(r.URL.Query().Get("mine_id"), 10, 64)
			vehicleID, _ := strconv.ParseUint(r.URL.Query().Get("vehicle_id"), 10, 64)
			resp, _ := client.ListDevice(r.Context(), &devicev1.ListDeviceRequest{
				Page:      int32(page),
				PageSize:  int32(pageSize),
				MineId:    mineID,
				VehicleId: vehicleID,
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/devices",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req devicev1.CreateDeviceRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CreateDevice(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/devices/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.GetDevice(r.Context(), &devicev1.GetDeviceRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPut,
		Path:   "/api/v1/devices/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			body, _ := io.ReadAll(r.Body)
			var req devicev1.UpdateDeviceRequest
			json.Unmarshal(body, &req)
			req.Id = id
			resp, _ := client.UpdateDevice(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodDelete,
		Path:   "/api/v1/devices/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.DeleteDevice(r.Context(), &devicev1.DeleteDeviceRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
