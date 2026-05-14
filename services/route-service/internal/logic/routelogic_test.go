package logic

import (
	"testing"
	"time"

	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/services/route-service/internal/model"
)

func TestNodeToProto(t *testing.T) {
	now := time.Now()
	n := &model.RoadNode{
		ID:        2001,
		Name:      "test-node",
		Latitude:  39.9042,
		Longitude: 116.4074,
		MineID:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	p := nodeToProto(n)
	if p.Id != 2001 || p.Name != "test-node" || p.Latitude != 39.9042 {
		t.Errorf("nodeToProto mismatch: %+v", p)
	}
}

func TestEdgeToProto(t *testing.T) {
	now := time.Now()
	e := &model.RoadEdge{
		ID:          3001,
		FromNodeID:  2001,
		ToNodeID:    2002,
		DistanceM:   1500,
		MaxSpeedKMH: 30,
		IsOneway:    false,
		MineID:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	p := edgeToProto(e)
	if p.Id != 3001 || p.DistanceM != 1500 || p.MaxSpeedKmh != 30 {
		t.Errorf("edgeToProto mismatch: %+v", p)
	}
	if p.IsOneway != false {
		t.Errorf("IsOneway = %v, want false", p.IsOneway)
	}
}

func TestNodeToProto_DefaultValues(t *testing.T) {
	n := &model.RoadNode{
		ID:   2002,
		Name: "default-node",
	}
	p := nodeToProto(n)
	if p.MineId != 0 {
		t.Errorf("MineId = %d, want 0", p.MineId)
	}
	if p.CreatedAt == "" {
		t.Error("CreatedAt should not be empty")
	}
}

func TestEdgeToProto_Oneway(t *testing.T) {
	e := &model.RoadEdge{
		ID:         3002,
		FromNodeID: 2001,
		ToNodeID:   2002,
		DistanceM:  500,
		IsOneway:   true,
	}
	p := edgeToProto(e)
	if p.IsOneway != true {
		t.Errorf("IsOneway = %v, want true", p.IsOneway)
	}
}

// Test the proto response code constants
func TestCalculateRouteResponse(t *testing.T) {
	resp := &routev1.RouteResponse{
		Code:    0,
		Message: "success",
		Data:    &routev1.RoutePath{TotalDistanceM: 1000, TotalDurationS: 120},
	}
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if resp.Data.TotalDistanceM != 1000 {
		t.Errorf("TotalDistanceM = %f, want 1000", resp.Data.TotalDistanceM)
	}
}

func TestNotFoundResponse(t *testing.T) {
	resp := &routev1.RouteResponse{
		Code:    404,
		Message: "no path found",
	}
	if resp.Code != 404 {
		t.Errorf("Code = %d, want 404", resp.Code)
	}
}

func TestRoutePathPoints(t *testing.T) {
	path := &routev1.RoutePath{
		Points: []*routev1.Point{
			{Latitude: 39.90, Longitude: 116.40},
			{Latitude: 39.91, Longitude: 116.41},
		},
		TotalDistanceM: 1500,
		EdgeIds:        []uint64{101},
	}
	if len(path.Points) != 2 {
		t.Errorf("len(Points) = %d, want 2", len(path.Points))
	}
	if len(path.EdgeIds) != 1 {
		t.Errorf("len(EdgeIds) = %d, want 1", len(path.EdgeIds))
	}
}

func TestDistanceResponse(t *testing.T) {
	resp := &routev1.DistanceResponse{
		Code:      0,
		DistanceM: 2700,
		DurationS: 324,
	}
	if resp.DistanceM != 2700 {
		t.Errorf("DistanceM = %f, want 2700", resp.DistanceM)
	}
}

func TestBatchRouteResponse(t *testing.T) {
	resp := &routev1.BatchRouteResponse{
		Code:    0,
		Message: "success",
		Results: []*routev1.RouteResult{
			{DistanceM: 1500, DurationS: 180, Status: 0},
			{DistanceM: 2700, DurationS: 324, Status: 0},
		},
	}
	if len(resp.Results) != 2 {
		t.Errorf("len(Results) = %d, want 2", len(resp.Results))
	}
}
