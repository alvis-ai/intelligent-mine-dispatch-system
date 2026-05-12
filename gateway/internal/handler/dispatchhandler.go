package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aicong/mine-dispatch/proto/dispatch/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterDispatchRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := dispatchv1.NewDispatchServiceClient(conn)

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/dispatch/assign",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req dispatchv1.AssignTaskRequest
			json.Unmarshal(body, &req)
			resp, _ := client.AssignTask(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/dispatch/tasks",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			resp, _ := client.ListTask(r.Context(), &dispatchv1.ListTaskRequest{Page: 1, PageSize: 50})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
