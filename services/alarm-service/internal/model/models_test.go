package model

import (
	"testing"
)

func TestGeofence_TableName(t *testing.T) {
	if (Geofence{}).TableName() != "geofences" {
		t.Errorf("Geofence table name should be geofences")
	}
}

func TestAlarmRule_TableName(t *testing.T) {
	if (AlarmRule{}).TableName() != "alarm_rules" {
		t.Errorf("AlarmRule table name should be alarm_rules")
	}
}

func TestAlarmEvent_TableName(t *testing.T) {
	if (AlarmEvent{}).TableName() != "alarm_events" {
		t.Errorf("AlarmEvent table name should be alarm_events")
	}
}

func TestSeverityConstants(t *testing.T) {
	tests := []struct {
		name, val string
	}{
		{"SeverityInfo", SeverityInfo},
		{"SeverityWarning", SeverityWarning},
		{"SeverityCritical", SeverityCritical},
	}
	for _, tt := range tests {
		if tt.val == "" {
			t.Errorf("%s should not be empty", tt.name)
		}
	}
}

func TestRuleTypeConstants(t *testing.T) {
	tests := []struct {
		name, val string
	}{
		{"RuleTypeGeofence", RuleTypeGeofence},
		{"RuleTypeSpeeding", RuleTypeSpeeding},
		{"RuleTypeOffline", RuleTypeOffline},
		{"RuleTypeDeviation", RuleTypeDeviation},
	}
	for _, tt := range tests {
		if tt.val == "" {
			t.Errorf("%s should not be empty", tt.name)
		}
	}
}

func TestShapeConstants(t *testing.T) {
	if ShapeCircle != "circle" {
		t.Errorf("ShapeCircle = %s, want circle", ShapeCircle)
	}
	if ShapePolygon != "polygon" {
		t.Errorf("ShapePolygon = %s, want polygon", ShapePolygon)
	}
}

func TestGeofenceDefaultValues(t *testing.T) {
	g := Geofence{}
	if g.Shape != "" {
		t.Errorf("default Shape should be empty, got %s", g.Shape)
	}
	if g.Enabled != false {
		t.Errorf("default Enabled should be false")
	}
}

func TestAlarmEventDefaultValues(t *testing.T) {
	e := AlarmEvent{}
	if e.Acknowledged != false {
		t.Errorf("default Acknowledged should be false")
	}
	if e.Severity != "" {
		t.Errorf("default Severity should be empty, got %s", e.Severity)
	}
}
