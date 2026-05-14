package server

import (
	"context"

	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/services/route-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/route-service/internal/svc"
)

type RouteServer struct {
	svc *svc.ServiceContext
	routev1.UnimplementedRouteServiceServer
}

func NewRouteServer(svc *svc.ServiceContext) *RouteServer {
	return &RouteServer{svc: svc}
}

func (s *RouteServer) CreateNode(ctx context.Context, in *routev1.CreateNodeRequest) (*routev1.NodeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).CreateNode(in)
}
func (s *RouteServer) UpdateNode(ctx context.Context, in *routev1.UpdateNodeRequest) (*routev1.NodeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).UpdateNode(in)
}
func (s *RouteServer) GetNode(ctx context.Context, in *routev1.GetNodeRequest) (*routev1.NodeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).GetNode(in)
}
func (s *RouteServer) DeleteNode(ctx context.Context, in *routev1.DeleteNodeRequest) (*routev1.NodeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).DeleteNode(in)
}
func (s *RouteServer) ListNodes(ctx context.Context, in *routev1.ListNodeRequest) (*routev1.NodeListResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).ListNodes(in)
}

func (s *RouteServer) CreateEdge(ctx context.Context, in *routev1.CreateEdgeRequest) (*routev1.EdgeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).CreateEdge(in)
}
func (s *RouteServer) UpdateEdge(ctx context.Context, in *routev1.UpdateEdgeRequest) (*routev1.EdgeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).UpdateEdge(in)
}
func (s *RouteServer) GetEdge(ctx context.Context, in *routev1.GetEdgeRequest) (*routev1.EdgeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).GetEdge(in)
}
func (s *RouteServer) DeleteEdge(ctx context.Context, in *routev1.DeleteEdgeRequest) (*routev1.EdgeResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).DeleteEdge(in)
}
func (s *RouteServer) ListEdges(ctx context.Context, in *routev1.ListEdgeRequest) (*routev1.EdgeListResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).ListEdges(in)
}

func (s *RouteServer) CalculateRoute(ctx context.Context, in *routev1.CalculateRouteRequest) (*routev1.RouteResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).CalculateRoute(in)
}
func (s *RouteServer) GetDistance(ctx context.Context, in *routev1.GetDistanceRequest) (*routev1.DistanceResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).GetDistance(in)
}
func (s *RouteServer) BatchCalculate(ctx context.Context, in *routev1.BatchRouteRequest) (*routev1.BatchRouteResponse, error) {
	return logic.NewRouteLogic(ctx, s.svc).BatchCalculate(in)
}
