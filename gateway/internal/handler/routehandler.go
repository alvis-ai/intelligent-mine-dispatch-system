package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterRouteRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := routev1.NewRouteServiceClient(conn)

	// ── Nodes ──
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/road-nodes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			mineID, _ := strconv.ParseUint(r.URL.Query().Get("mine_id"), 10, 64)
			resp, _ := client.ListNodes(r.Context(), &routev1.ListNodeRequest{MineId: mineID})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/road-nodes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req routev1.CreateNodeRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CreateNode(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/road-nodes/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.GetNode(r.Context(), &routev1.GetNodeRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPut,
		Path:   "/api/v1/road-nodes/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			body, _ := io.ReadAll(r.Body)
			var req routev1.UpdateNodeRequest
			json.Unmarshal(body, &req)
			req.Id = id
			resp, _ := client.UpdateNode(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodDelete,
		Path:   "/api/v1/road-nodes/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.DeleteNode(r.Context(), &routev1.DeleteNodeRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	// ── Edges ──
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/road-edges",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			mineID, _ := strconv.ParseUint(r.URL.Query().Get("mine_id"), 10, 64)
			nodeID, _ := strconv.ParseUint(r.URL.Query().Get("node_id"), 10, 64)
			resp, _ := client.ListEdges(r.Context(), &routev1.ListEdgeRequest{MineId: mineID, NodeId: nodeID})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/road-edges",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req routev1.CreateEdgeRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CreateEdge(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPut,
		Path:   "/api/v1/road-edges/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			body, _ := io.ReadAll(r.Body)
			var req routev1.UpdateEdgeRequest
			json.Unmarshal(body, &req)
			req.Id = id
			resp, _ := client.UpdateEdge(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodDelete,
		Path:   "/api/v1/road-edges/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.DeleteEdge(r.Context(), &routev1.DeleteEdgeRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	// ── Routing ──
	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/route/calculate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req routev1.CalculateRouteRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CalculateRoute(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/route/distance",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req routev1.GetDistanceRequest
			json.Unmarshal(body, &req)
			resp, _ := client.GetDistance(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/route/batch",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req routev1.BatchRouteRequest
			json.Unmarshal(body, &req)
			resp, _ := client.BatchCalculate(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}
