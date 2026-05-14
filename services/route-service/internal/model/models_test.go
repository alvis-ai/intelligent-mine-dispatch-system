package model

import "testing"

func TestRoadNode_TableName(t *testing.T) {
	if (RoadNode{}).TableName() != "road_nodes" {
		t.Errorf("RoadNode table name should be road_nodes")
	}
}

func TestRoadEdge_TableName(t *testing.T) {
	if (RoadEdge{}).TableName() != "road_edges" {
		t.Errorf("RoadEdge table name should be road_edges")
	}
}

func TestRoadNode_Defaults(t *testing.T) {
	n := RoadNode{ID: 1, Name: "test", Latitude: 39.9, Longitude: 116.4}
	if n.MineID != 0 {
		t.Errorf("default MineID = %d, want 0", n.MineID)
	}
}

func TestRoadEdge_Defaults(t *testing.T) {
	e := RoadEdge{ID: 1, FromNodeID: 1, ToNodeID: 2, DistanceM: 1000}
	if e.MaxSpeedKMH != 0 {
		t.Errorf("default MaxSpeedKMH = %d, want 0", e.MaxSpeedKMH)
	}
	if e.IsOneway != false {
		t.Errorf("default IsOneway = %v, want false", e.IsOneway)
	}
}

func TestRoadEdge_Oneway(t *testing.T) {
	e := RoadEdge{
		ID: 1, FromNodeID: 1, ToNodeID: 2,
		DistanceM: 500, MaxSpeedKMH: 30, IsOneway: true,
	}
	if !e.IsOneway {
		t.Error("IsOneway should be true")
	}
}
