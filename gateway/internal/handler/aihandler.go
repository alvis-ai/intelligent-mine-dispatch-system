package handler

import (
	"encoding/json"
	"io"
	"net/http"

	aiv1 "github.com/aicong/mine-dispatch/proto/ai/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterAiRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := aiv1.NewAiServiceClient(conn)

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/ai/congestion",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req aiv1.PredictCongestionRequest
			json.Unmarshal(body, &req)
			resp, _ := client.PredictCongestion(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/ai/route/recommend",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req aiv1.RecommendRouteRequest
			json.Unmarshal(body, &req)
			resp, _ := client.RecommendRoute(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/ai/demand",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req aiv1.PredictDemandRequest
			json.Unmarshal(body, &req)
			resp, _ := client.PredictDemand(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/ai/suggest-assign",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req aiv1.SuggestAssignRequest
			json.Unmarshal(body, &req)
			resp, _ := client.SuggestAssign(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
