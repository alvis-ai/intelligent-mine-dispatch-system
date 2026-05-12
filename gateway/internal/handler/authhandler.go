package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterAuthRoutes(server *rest.Server, userRpc userv1.UserServiceClient) {
	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/users",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req userv1.CreateUserRequest
			json.Unmarshal(body, &req)
			resp, _ := userRpc.CreateUser(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/users",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			page := 1
			pageSize := 20
			resp, _ := userRpc.ListUser(r.Context(), &userv1.ListUserRequest{
				Page:     int32(page),
				PageSize: int32(pageSize),
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
