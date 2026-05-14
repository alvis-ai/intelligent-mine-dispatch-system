package logic

import (
	"context"
	"time"

	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/services/route-service/internal/model"
	"github.com/aicong/mine-dispatch/services/route-service/internal/svc"
	"github.com/aicong/mine-dispatch/pkg/utils"
)

type RouteLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewRouteLogic(ctx context.Context, svc *svc.ServiceContext) *RouteLogic {
	return &RouteLogic{ctx: ctx, svc: svc}
}

// ── Node CRUD ──

func (l *RouteLogic) CreateNode(in *routev1.CreateNodeRequest) (*routev1.NodeResponse, error) {
	n := model.RoadNode{
		ID:        utils.NextID(),
		Name:      in.Name,
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
		MineID:    in.MineId,
	}
	if err := l.svc.DB.Create(&n).Error; err != nil {
		return &routev1.NodeResponse{Code: 500, Message: err.Error()}, nil
	}
	return &routev1.NodeResponse{Code: 0, Message: "success", Data: nodeToProto(&n)}, nil
}

func (l *RouteLogic) UpdateNode(in *routev1.UpdateNodeRequest) (*routev1.NodeResponse, error) {
	var n model.RoadNode
	if err := l.svc.DB.First(&n, in.Id).Error; err != nil {
		return &routev1.NodeResponse{Code: 404, Message: "node not found"}, nil
	}
	updates := map[string]interface{}{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Latitude != 0 {
		updates["latitude"] = in.Latitude
	}
	if in.Longitude != 0 {
		updates["longitude"] = in.Longitude
	}
	if in.MineId != 0 {
		updates["mine_id"] = in.MineId
	}
	l.svc.DB.Model(&n).Updates(updates)
	l.svc.DB.First(&n, in.Id)
	return &routev1.NodeResponse{Code: 0, Message: "success", Data: nodeToProto(&n)}, nil
}

func (l *RouteLogic) GetNode(in *routev1.GetNodeRequest) (*routev1.NodeResponse, error) {
	var n model.RoadNode
	if err := l.svc.DB.First(&n, in.Id).Error; err != nil {
		return &routev1.NodeResponse{Code: 404, Message: "node not found"}, nil
	}
	return &routev1.NodeResponse{Code: 0, Message: "success", Data: nodeToProto(&n)}, nil
}

func (l *RouteLogic) DeleteNode(in *routev1.DeleteNodeRequest) (*routev1.NodeResponse, error) {
	l.svc.DB.Delete(&model.RoadNode{}, in.Id)
	l.svc.DB.Where("from_node_id = ? OR to_node_id = ?", in.Id, in.Id).Delete(&model.RoadEdge{})
	return &routev1.NodeResponse{Code: 0, Message: "deleted"}, nil
}

func (l *RouteLogic) ListNodes(in *routev1.ListNodeRequest) (*routev1.NodeListResponse, error) {
	var nodes []model.RoadNode
	var total int64
	db := l.svc.DB.Model(&model.RoadNode{})
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	db.Count(&total)
	db.Order("id ASC").Find(&nodes)
	var list []*routev1.RoadNode
	for i := range nodes {
		list = append(list, nodeToProto(&nodes[i]))
	}
	return &routev1.NodeListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}

// ── Edge CRUD ──

func (l *RouteLogic) CreateEdge(in *routev1.CreateEdgeRequest) (*routev1.EdgeResponse, error) {
	e := model.RoadEdge{
		ID:          utils.NextID(),
		FromNodeID:  in.FromNodeId,
		ToNodeID:    in.ToNodeId,
		DistanceM:   in.DistanceM,
		MaxSpeedKMH: in.MaxSpeedKmh,
		IsOneway:    in.IsOneway,
		MineID:      in.MineId,
	}
	if err := l.svc.DB.Create(&e).Error; err != nil {
		return &routev1.EdgeResponse{Code: 500, Message: err.Error()}, nil
	}
	return &routev1.EdgeResponse{Code: 0, Message: "success", Data: edgeToProto(&e)}, nil
}

func (l *RouteLogic) UpdateEdge(in *routev1.UpdateEdgeRequest) (*routev1.EdgeResponse, error) {
	var e model.RoadEdge
	if err := l.svc.DB.First(&e, in.Id).Error; err != nil {
		return &routev1.EdgeResponse{Code: 404, Message: "edge not found"}, nil
	}
	updates := map[string]interface{}{}
	if in.FromNodeId != 0 {
		updates["from_node_id"] = in.FromNodeId
	}
	if in.ToNodeId != 0 {
		updates["to_node_id"] = in.ToNodeId
	}
	if in.DistanceM != 0 {
		updates["distance_m"] = in.DistanceM
	}
	if in.MaxSpeedKmh != 0 {
		updates["max_speed_kmh"] = in.MaxSpeedKmh
	}
	updates["is_oneway"] = in.IsOneway
	if in.MineId != 0 {
		updates["mine_id"] = in.MineId
	}
	l.svc.DB.Model(&e).Updates(updates)
	l.svc.DB.First(&e, in.Id)
	return &routev1.EdgeResponse{Code: 0, Message: "success", Data: edgeToProto(&e)}, nil
}

func (l *RouteLogic) GetEdge(in *routev1.GetEdgeRequest) (*routev1.EdgeResponse, error) {
	var e model.RoadEdge
	if err := l.svc.DB.First(&e, in.Id).Error; err != nil {
		return &routev1.EdgeResponse{Code: 404, Message: "edge not found"}, nil
	}
	return &routev1.EdgeResponse{Code: 0, Message: "success", Data: edgeToProto(&e)}, nil
}

func (l *RouteLogic) DeleteEdge(in *routev1.DeleteEdgeRequest) (*routev1.EdgeResponse, error) {
	l.svc.DB.Delete(&model.RoadEdge{}, in.Id)
	return &routev1.EdgeResponse{Code: 0, Message: "deleted"}, nil
}

func (l *RouteLogic) ListEdges(in *routev1.ListEdgeRequest) (*routev1.EdgeListResponse, error) {
	var edges []model.RoadEdge
	var total int64
	db := l.svc.DB.Model(&model.RoadEdge{})
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	if in.NodeId > 0 {
		db = db.Where("from_node_id = ? OR to_node_id = ?", in.NodeId, in.NodeId)
	}
	db.Count(&total)
	db.Order("id ASC").Find(&edges)
	var list []*routev1.RoadEdge
	for i := range edges {
		list = append(list, edgeToProto(&edges[i]))
	}
	return &routev1.EdgeListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}

// ── Routing ──

func (l *RouteLogic) CalculateRoute(in *routev1.CalculateRouteRequest) (*routev1.RouteResponse, error) {
	var nodes []model.RoadNode
	var edges []model.RoadEdge
	db := l.svc.DB.Model(&model.RoadNode{})
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	db.Find(&nodes)
	edgeDB := l.svc.DB.Model(&model.RoadEdge{})
	if in.MineId > 0 {
		edgeDB = edgeDB.Where("mine_id = ?", in.MineId)
	}
	edgeDB.Find(&edges)

	g := BuildGraph(nodes, edges)
	fromID, _ := g.FindNearestNode(in.FromLat, in.FromLon)
	toID, _ := g.FindNearestNode(in.ToLat, in.ToLon)

	if fromID == 0 || toID == 0 {
		return &routev1.RouteResponse{Code: 404, Message: "no nearby road nodes found"}, nil
	}

	var result *PathResult
	switch in.Algorithm {
	case "astar":
		result = g.AStar(fromID, toID)
	default:
		result = g.Dijkstra(fromID, toID)
	}

	if result == nil || !result.Found {
		return &routev1.RouteResponse{Code: 404, Message: "no path found"}, nil
	}

	return &routev1.RouteResponse{
		Code:    0,
		Message: "success",
		Data:    g.PathToProto(result),
	}, nil
}

func (l *RouteLogic) GetDistance(in *routev1.GetDistanceRequest) (*routev1.DistanceResponse, error) {
	// For a quick distance estimate, use Haversine as fallback
	directDist := haversine(in.FromLat, in.FromLon, in.ToLat, in.ToLon)

	// Try to find road distance
	var nodes []model.RoadNode
	var edges []model.RoadEdge
	db := l.svc.DB.Model(&model.RoadNode{})
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	db.Find(&nodes)
	l.svc.DB.Model(&model.RoadEdge{}).Find(&edges)

	g := BuildGraph(nodes, edges)
	fromID, _ := g.FindNearestNode(in.FromLat, in.FromLon)
	toID, _ := g.FindNearestNode(in.ToLat, in.ToLon)

	if fromID == 0 || toID == 0 || fromID == toID {
		dur := directDist / 1000 / 30 * 3600
		return &routev1.DistanceResponse{
			Code:       0,
			Message:    "success",
			DistanceM:  directDist,
			DurationS:  dur,
		}, nil
	}

	result := g.Dijkstra(fromID, toID)
	if result != nil && result.Found {
		return &routev1.DistanceResponse{
			Code:       0,
			Message:    "success",
			DistanceM:  result.TotalDistM,
			DurationS:  result.TotalDurS,
		}, nil
	}

	dur := directDist / 1000 / 30 * 3600
	return &routev1.DistanceResponse{
		Code:       0,
		Message:    "success",
		DistanceM:  directDist,
		DurationS:  dur,
	}, nil
}

func (l *RouteLogic) BatchCalculate(in *routev1.BatchRouteRequest) (*routev1.BatchRouteResponse, error) {
	var results []*routev1.RouteResult
	for _, req := range in.Requests {
		distResp, err := l.GetDistance(&routev1.GetDistanceRequest{
			FromLat:  req.FromLat,
			FromLon:  req.FromLon,
			ToLat:    req.ToLat,
			ToLon:    req.ToLon,
			MineId:   req.MineId,
		})
		if err != nil {
			results = append(results, &routev1.RouteResult{Status: 1, Error: err.Error()})
		} else {
			results = append(results, &routev1.RouteResult{
				DistanceM: distResp.DistanceM,
				DurationS: distResp.DurationS,
				Status:    0,
			})
		}
	}
	return &routev1.BatchRouteResponse{Code: 0, Message: "success", Results: results}, nil
}

// ── Conversions ──

func nodeToProto(n *model.RoadNode) *routev1.RoadNode {
	return &routev1.RoadNode{
		Id:        n.ID,
		Name:      n.Name,
		Latitude:  n.Latitude,
		Longitude: n.Longitude,
		MineId:    n.MineID,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
		UpdatedAt: n.UpdatedAt.Format(time.RFC3339),
	}
}

func edgeToProto(e *model.RoadEdge) *routev1.RoadEdge {
	return &routev1.RoadEdge{
		Id:          e.ID,
		FromNodeId:  e.FromNodeID,
		ToNodeId:    e.ToNodeID,
		DistanceM:   e.DistanceM,
		MaxSpeedKmh: e.MaxSpeedKMH,
		IsOneway:    e.IsOneway,
		MineId:      e.MineID,
		CreatedAt:   e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   e.UpdatedAt.Format(time.RFC3339),
	}
}
