package logic

import (
	"testing"
	"time"

	alarmv1 "github.com/aicong/mine-dispatch/proto/alarm/v1"
	"github.com/aicong/mine-dispatch/services/alarm-service/internal/model"
)

func TestModelToProtoGeofence(t *testing.T) {
	now := time.Now()
	g := &model.Geofence{
		ID:          1001,
		Name:        "test-fence",
		Shape:       "circle",
		CenterLat:   39.9042,
		CenterLon:   116.4074,
		RadiusM:     500,
		PointsJSON:  `[{"lat":39.9,"lon":116.4},{"lat":39.91,"lon":116.41}]`,
		FenceType:   "restricted",
		MinSpeedKMH: 0,
		MaxSpeedKMH: 40,
		TimeRange:   "00:00-23:59",
		Enabled:     true,
		MineID:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	p := modelToProtoGeofence(g)
	if p.Id != 1001 {
		t.Errorf("Id = %d, want 1001", p.Id)
	}
	if p.Name != "test-fence" {
		t.Errorf("Name = %s, want test-fence", p.Name)
	}
	if p.Shape != "circle" {
		t.Errorf("Shape = %s, want circle", p.Shape)
	}
	if p.CenterLat != 39.9042 {
		t.Errorf("CenterLat = %f, want 39.9042", p.CenterLat)
	}
	if p.RadiusM != 500 {
		t.Errorf("RadiusM = %f, want 500", p.RadiusM)
	}
	if len(p.Points) != 2 {
		t.Errorf("len(Points) = %d, want 2", len(p.Points))
	}
	if p.Enabled != true {
		t.Errorf("Enabled = %v, want true", p.Enabled)
	}
}

func TestModelToProtoGeofence_EmptyPoints(t *testing.T) {
	g := &model.Geofence{
		ID:        1002,
		Name:      "no-points",
		Shape:     "circle",
		RadiusM:   100,
		FenceType: "loading",
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	p := modelToProtoGeofence(g)
	if p == nil {
		t.Fatal("modelToProtoGeofence returned nil")
	}
	if p.RadiusM != 100 {
		t.Errorf("RadiusM = %f, want 100", p.RadiusM)
	}
}

func TestModelToProtoRule(t *testing.T) {
	now := time.Now()
	r := &model.AlarmRule{
		ID:          5001,
		Name:        "speeding-rule",
		RuleType:    "speeding",
		GeofenceID:  1001,
		Severity:    "critical",
		Description: "speed > 80 km/h",
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	p := modelToProtoRule(r)
	if p.Id != 5001 {
		t.Errorf("Id = %d, want 5001", p.Id)
	}
	if p.Name != "speeding-rule" {
		t.Errorf("Name = %s", p.Name)
	}
	if p.RuleType != "speeding" {
		t.Errorf("RuleType = %s", p.RuleType)
	}
	if p.Severity != "critical" {
		t.Errorf("Severity = %s", p.Severity)
	}
	if p.Enabled != true {
		t.Errorf("Enabled = %v, want true", p.Enabled)
	}
}

func TestModelToProtoEvent(t *testing.T) {
	now := time.Now()
	ackAt := now.Add(-1 * time.Hour)

	tests := []struct {
		name string
		event *model.AlarmEvent
		check func(*testing.T, *alarmv1.AlarmEvent)
	}{
		{
			name: "unacknowledged event",
			event: &model.AlarmEvent{
				ID:       9001,
				RuleID:   5001,
				VehicleID: 101,
				VehiclePlate: "矿卡-A001",
				AlarmType: "speeding",
				Severity:  "critical",
				Message:   "超速: 85 km/h",
				Latitude:  39.9042,
				Longitude: 116.4074,
				Speed:     85,
				MineID:    1,
				CreatedAt: now,
			},
			check: func(t *testing.T, p *alarmv1.AlarmEvent) {
				if p.Id != 9001 { t.Errorf("Id = %d, want 9001", p.Id) }
				if p.AlarmType != "speeding" { t.Errorf("AlarmType = %s", p.AlarmType) }
				if p.Acknowledged != false { t.Errorf("Acknowledged = %v, want false", p.Acknowledged) }
				if p.AcknowledgedAt != "" { t.Errorf("AcknowledgedAt = %s, want empty", p.AcknowledgedAt) }
				if p.Speed != 85 { t.Errorf("Speed = %f, want 85", p.Speed) }
			},
		},
		{
			name: "acknowledged event",
			event: &model.AlarmEvent{
				ID:             9002,
				RuleID:         5001,
				VehicleID:      102,
				AlarmType:      "geofence_entry",
				Severity:       "warning",
				Acknowledged:   true,
				AcknowledgedBy: "admin",
				AcknowledgedAt: &ackAt,
				CreatedAt:      now,
			},
			check: func(t *testing.T, p *alarmv1.AlarmEvent) {
				if p.Acknowledged != true { t.Errorf("Acknowledged = %v, want true", p.Acknowledged) }
				if p.AcknowledgedBy != "admin" { t.Errorf("AcknowledgedBy = %s, want admin", p.AcknowledgedBy) }
				if p.AcknowledgedAt == "" { t.Errorf("AcknowledgedAt should not be empty") }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := modelToProtoEvent(tt.event)
			if p == nil {
				t.Fatal("modelToProtoEvent returned nil")
			}
			tt.check(t, p)
		})
	}
}

func TestPointInCircle_EdgeCases(t *testing.T) {
	// Point exactly at center
	if !pointInCircle(39.9, 116.4, 39.9, 116.4, 100) {
		t.Error("center point should be inside")
	}
	// Zero radius — only center matches
	if pointInCircle(39.901, 116.4, 39.9, 116.4, 0) {
		t.Error("point at distance should be outside when radius is 0")
	}
	// Negative radius — point at center has distance=0, and 0 <= -1 is false
	if pointInCircle(39.9, 116.4, 39.9, 116.4, -1) {
		t.Error("point should not be inside with negative radius")
	}
}

func TestModelToProtoGeofence_DefaultValues(t *testing.T) {
	g := &model.Geofence{
		ID:        999,
		Name:      "defaults",
		Shape:     "polygon",
		FenceType: "dumping",
		Enabled:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	p := modelToProtoGeofence(g)
	if p.Name != "defaults" {
		t.Errorf("Name = %s, want defaults", p.Name)
	}
	if p.Enabled != false {
		t.Errorf("Enabled = %v, want false", p.Enabled)
	}
	if p.FenceType != "dumping" {
		t.Errorf("FenceType = %s, want dumping", p.FenceType)
	}
}
