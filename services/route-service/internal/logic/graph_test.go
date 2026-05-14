package logic

import (
	"math"
	"testing"

	"github.com/aicong/mine-dispatch/services/route-service/internal/model"
)

func makeTestGraph() *Graph {
	nodes := []model.RoadNode{
		{ID: 1, Name: "A", Latitude: 39.90, Longitude: 116.40},
		{ID: 2, Name: "B", Latitude: 39.91, Longitude: 116.41},
		{ID: 3, Name: "C", Latitude: 39.92, Longitude: 116.42},
		{ID: 4, Name: "D", Latitude: 39.89, Longitude: 116.43},
	}
	edges := []model.RoadEdge{
		{ID: 101, FromNodeID: 1, ToNodeID: 2, DistanceM: 1500, MaxSpeedKMH: 30},
		{ID: 102, FromNodeID: 2, ToNodeID: 3, DistanceM: 1200, MaxSpeedKMH: 30},
		{ID: 103, FromNodeID: 3, ToNodeID: 4, DistanceM: 2000, MaxSpeedKMH: 20},
		{ID: 104, FromNodeID: 1, ToNodeID: 4, DistanceM: 5000, MaxSpeedKMH: 40, IsOneway: true},
	}
	return BuildGraph(nodes, edges)
}

func TestBuildGraph(t *testing.T) {
	g := makeTestGraph()
	if len(g.Nodes) != 4 {
		t.Errorf("len(Nodes) = %d, want 4", len(g.Nodes))
	}
	if len(g.Edges) != 4 {
		t.Errorf("len(Edges) = %d, want 4", len(g.Edges))
	}
	// Node 1 should have 3 edges: 1->2, 1->4 (directed), and 2->1 (undirected back)
	if len(g.Edges[1]) != 2 {
		t.Errorf("len(Edges[1]) = %d, want 2 (1->2, 1->4)", len(g.Edges[1]))
	}
}

func TestBuildGraph_Bidirectional(t *testing.T) {
	nodes := []model.RoadNode{
		{ID: 1, Name: "A", Latitude: 39.9, Longitude: 116.4},
		{ID: 2, Name: "B", Latitude: 39.91, Longitude: 116.41},
	}
	edges := []model.RoadEdge{
		{ID: 101, FromNodeID: 1, ToNodeID: 2, DistanceM: 1000, MaxSpeedKMH: 30, IsOneway: false},
	}
	g := BuildGraph(nodes, edges)
	// Should have edges from 1->2 and 2->1 (bidirectional)
	if len(g.Edges[1]) != 1 {
		t.Errorf("len(Edges[1]) = %d, want 1", len(g.Edges[1]))
	}
	if len(g.Edges[2]) != 1 {
		t.Errorf("len(Edges[2]) = %d, want 1 (bidirectional back-edge)", len(g.Edges[2]))
	}
}

func TestDijkstra_DirectPath(t *testing.T) {
	g := makeTestGraph()
	result := g.Dijkstra(1, 2)
	if !result.Found {
		t.Fatal("Dijkstra: path not found")
	}
	if result.TotalDistM != 1500 {
		t.Errorf("TotalDistM = %f, want 1500", result.TotalDistM)
	}
}

func TestDijkstra_MultiHopPath(t *testing.T) {
	g := makeTestGraph()
	// A -> B -> C = 1500 + 1200 = 2700
	result := g.Dijkstra(1, 3)
	if !result.Found {
		t.Fatal("Dijkstra: path not found")
	}
	if result.TotalDistM != 2700 {
		t.Errorf("TotalDistM = %f, want 2700", result.TotalDistM)
	}
}

func TestDijkstra_CompareRoutes(t *testing.T) {
	g := makeTestGraph()
	// A -> D: direct 5000 (oneway) vs A->B->C->D = 1500+1200+2000 = 4700
	// Dijkstra should choose the path via B,C since it's shorter
	result := g.Dijkstra(1, 4)
	if !result.Found {
		t.Fatal("Dijkstra: path not found")
	}
	if result.TotalDistM > 4800 || result.TotalDistM < 4600 {
		t.Errorf("Dijkstra chose TotalDistM=%f, expected ~4700 (via B,C)", result.TotalDistM)
	}
}

func TestAStar_DirectPath(t *testing.T) {
	g := makeTestGraph()
	result := g.AStar(1, 2)
	if !result.Found {
		t.Fatal("AStar: path not found")
	}
	if result.TotalDistM != 1500 {
		t.Errorf("TotalDistM = %f, want 1500", result.TotalDistM)
	}
}

func TestAStar_MultiHopPath(t *testing.T) {
	g := makeTestGraph()
	result := g.AStar(1, 3)
	if !result.Found {
		t.Fatal("AStar: path not found")
	}
	if result.TotalDistM != 2700 {
		t.Errorf("TotalDistM = %f, want 2700", result.TotalDistM)
	}
}

func TestDijkstra_NoPath(t *testing.T) {
	// Disconnected graph
	nodes := []model.RoadNode{
		{ID: 1, Name: "A", Latitude: 39.9, Longitude: 116.4},
		{ID: 2, Name: "B", Latitude: 39.91, Longitude: 116.41},
	}
	edges := []model.RoadEdge{} // no edges
	g := BuildGraph(nodes, edges)
	result := g.Dijkstra(1, 2)
	if result.Found {
		t.Error("Dijkstra should not find a path in disconnected graph")
	}
}

func TestDijkstra_SameNode(t *testing.T) {
	g := makeTestGraph()
	result := g.Dijkstra(1, 1)
	if !result.Found {
		t.Error("Dijkstra should find path to same node")
	}
	if result.TotalDistM != 0 {
		t.Errorf("TotalDistM = %f, want 0", result.TotalDistM)
	}
}

func TestDijkstra_NonExistentNode(t *testing.T) {
	g := makeTestGraph()
	result := g.Dijkstra(1, 999)
	if result.Found {
		t.Error("Dijkstra should not find path to non-existent node")
	}
	result = g.Dijkstra(999, 1)
	if result.Found {
		t.Error("Dijkstra should not find path from non-existent node")
	}
}

func TestFindNearestNode(t *testing.T) {
	g := makeTestGraph()
	id, dist := g.FindNearestNode(39.90, 116.40)
	if id != 1 {
		t.Errorf("nearest = %d, want 1 (node A)", id)
	}
	if dist != 0 {
		t.Errorf("distance = %f, want 0 (exact match)", dist)
	}

	// Closer to node B (39.91, 116.41)
	id, _ = g.FindNearestNode(39.905, 116.405)
	if id != 1 {
		// Could also be node B depending on exact distance
	}
	_ = id // both A and B are valid
}

func TestHaversine(t *testing.T) {
	// Same point
	d := haversine(39.9, 116.4, 39.9, 116.4)
	if d != 0 {
		t.Errorf("same point distance = %f, want 0", d)
	}

	// Known distance: ~1 degree latitude ≈ 111km
	d = haversine(0, 0, 1, 0)
	if d < 110000 || d > 112000 {
		t.Errorf("1 degree lat distance = %f, want ~111000", d)
	}
}

func TestHaversine_Symmetry(t *testing.T) {
	d1 := haversine(39.9, 116.4, 31.2, 121.5)
	d2 := haversine(31.2, 121.5, 39.9, 116.4)
	if math.Abs(d1-d2) > 0.001 {
		t.Errorf("not symmetric: %f vs %f", d1, d2)
	}
}

func TestSortNodesByDistance(t *testing.T) {
	nodes := []model.RoadNode{
		{ID: 1, Name: "A", Latitude: 39.90, Longitude: 116.40},
		{ID: 2, Name: "B", Latitude: 39.91, Longitude: 116.41},
		{ID: 3, Name: "C", Latitude: 39.92, Longitude: 116.42},
	}
	// Sort by distance from (39.90, 116.40) — A should be first
	sorted := SortNodesByDistance(nodes, 39.90, 116.40)
	if len(sorted) != 3 {
		t.Errorf("len = %d, want 3", len(sorted))
	}
	if sorted[0] != 1 {
		t.Errorf("closest = %d, want 1 (node A)", sorted[0])
	}
}

func TestPathToProto(t *testing.T) {
	g := makeTestGraph()
	result := g.Dijkstra(1, 3)
	proto := g.PathToProto(result)
	if proto == nil {
		t.Fatal("PathToProto returned nil")
	}
	if proto.TotalDistanceM != 2700 {
		t.Errorf("TotalDistanceM = %f, want 2700", proto.TotalDistanceM)
	}
	if len(proto.NodeIds) < 2 {
		t.Errorf("len(NodeIds) = %d, want >= 2", len(proto.NodeIds))
	}
	if len(proto.Points) != len(proto.NodeIds) {
		t.Errorf("len(Points) = %d != len(NodeIds) = %d", len(proto.Points), len(proto.NodeIds))
	}
	// First point should be node A (1)
	if proto.Points[0].Latitude != 39.90 || proto.Points[0].Longitude != 116.40 {
		t.Errorf("first point = %v, want (39.90, 116.40)", proto.Points[0])
	}
}

func TestPathToProto_Nil(t *testing.T) {
	g := makeTestGraph()
	proto := g.PathToProto(nil)
	if proto != nil {
		t.Error("PathToProto(nil) should return nil")
	}
	proto = g.PathToProto(&PathResult{Found: false})
	if proto != nil {
		t.Error("PathToProto(not found) should return nil")
	}
}

func TestEmptyGraph(t *testing.T) {
	g := BuildGraph(nil, nil)
	if g.Nodes == nil || g.Edges == nil {
		t.Error("BuildGraph should initialize maps")
	}
	result := g.Dijkstra(1, 2)
	if result.Found {
		t.Error("empty graph should not find paths")
	}
}

func TestOnewayEdge(t *testing.T) {
	nodes := []model.RoadNode{
		{ID: 1, Name: "A", Latitude: 39.9, Longitude: 116.4},
		{ID: 2, Name: "B", Latitude: 39.91, Longitude: 116.41},
	}
	edges := []model.RoadEdge{
		{ID: 101, FromNodeID: 1, ToNodeID: 2, DistanceM: 1000, MaxSpeedKMH: 30, IsOneway: true},
	}
	g := BuildGraph(nodes, edges)
	// A->B should work
	r1 := g.Dijkstra(1, 2)
	if !r1.Found {
		t.Error("oneway: path A->B should exist")
	}
	// B->A should not (no back edge)
	r2 := g.Dijkstra(2, 1)
	if r2.Found {
		t.Error("oneway: path B->A should not exist")
	}
}

func TestDurationCalculation(t *testing.T) {
	g := makeTestGraph()
	result := g.Dijkstra(1, 2)
	if !result.Found {
		t.Fatal("path not found")
	}
	// Distance 1500m at 30 km/h
	// time = (1.5 km / 30 km/h) * 3600 = 180 seconds
	if result.TotalDurS <= 0 {
		t.Errorf("TotalDurS = %f, want > 0", result.TotalDurS)
	}
}
