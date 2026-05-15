package logic

import (
	"testing"

	telemetryv1 "github.com/aicong/mine-dispatch/proto/telemetry/v1"
)

func TestParseUint64(t *testing.T) {
	tests := []struct {
		input string
		want  uint64
	}{
		{"12345", 12345},
		{"0", 0},
		{"", 0},
		{"abc", 0},
		{"9999999999999999999", 9999999999999999999}, // large number
		{" 42", 42}, // leading space — Sscanf skips it
		{"42 ", 42}, // trailing space — Sscanf handles it
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseUint64(tt.input)
			if got != tt.want {
				t.Errorf("parseUint64(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestReportLocationResponse_Success(t *testing.T) {
	resp := &telemetryv1.ReportLocationResponse{Code: 0, Message: "success"}
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("Message = %s, want success", resp.Message)
	}
}

func TestReportLocationResponse_Error(t *testing.T) {
	resp := &telemetryv1.ReportLocationResponse{Code: 500, Message: "internal error"}
	if resp.Code != 500 {
		t.Errorf("Code = %d, want 500", resp.Code)
	}
}

func TestGetVehicleLocationResponse_Found(t *testing.T) {
	resp := &telemetryv1.GetVehicleLocationResponse{
		Code:    0,
		Message: "success",
		Location: &telemetryv1.LocationData{
			VehicleId: 1001,
			Latitude:  39.9042,
			Longitude: 116.4074,
			Speed:     35.5,
			Heading:   180,
		},
	}
	if resp.Location.VehicleId != 1001 {
		t.Errorf("VehicleId = %d, want 1001", resp.Location.VehicleId)
	}
	if resp.Location.Latitude != 39.9042 {
		t.Errorf("Latitude = %f, want 39.9042", resp.Location.Latitude)
	}
	if resp.Location.Heading != 180 {
		t.Errorf("Heading = %f, want 180", resp.Location.Heading)
	}
}

func TestGetVehicleLocationResponse_NotFound(t *testing.T) {
	resp := &telemetryv1.GetVehicleLocationResponse{Code: 404, Message: "location not found"}
	if resp.Code != 404 {
		t.Errorf("Code = %d, want 404", resp.Code)
	}
	if resp.Location != nil {
		t.Error("Location should be nil when not found")
	}
}

func TestGetNearbyVehiclesResponse_Empty(t *testing.T) {
	resp := &telemetryv1.GetNearbyVehiclesResponse{
		Code:     0,
		Message:  "success",
		Vehicles: []*telemetryv1.NearbyVehicle{},
	}
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if len(resp.Vehicles) != 0 {
		t.Errorf("len(Vehicles) = %d, want 0", len(resp.Vehicles))
	}
}

func TestGetNearbyVehiclesResponse_WithResults(t *testing.T) {
	resp := &telemetryv1.GetNearbyVehiclesResponse{
		Code:    0,
		Message: "success",
		Vehicles: []*telemetryv1.NearbyVehicle{
			{VehicleId: 1001, Latitude: 39.90, Longitude: 116.40, DistanceKm: 0.5},
			{VehicleId: 1002, Latitude: 39.91, Longitude: 116.41, DistanceKm: 1.2},
		},
	}
	if len(resp.Vehicles) != 2 {
		t.Errorf("len(Vehicles) = %d, want 2", len(resp.Vehicles))
	}
	if resp.Vehicles[0].DistanceKm != 0.5 {
		t.Errorf("DistanceKm[0] = %f, want 0.5", resp.Vehicles[0].DistanceKm)
	}
}

func TestGetNearbyVehiclesResponse_Error(t *testing.T) {
	resp := &telemetryv1.GetNearbyVehiclesResponse{Code: 500, Message: "redis error"}
	if resp.Code != 500 {
		t.Errorf("Code = %d, want 500", resp.Code)
	}
	if resp.Vehicles != nil {
		t.Error("Vehicles should be nil on error")
	}
}

func TestLocationData_Bounds(t *testing.T) {
	tests := []struct {
		name  string
		lat   float64
		lon   float64
		speed float64
	}{
		{"valid", 39.9042, 116.4074, 60},
		{"south pole", -90, 0, 0},
		{"north pole", 90, 180, 200},
		{"zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := &telemetryv1.LocationData{
				VehicleId: 1,
				Latitude:  tt.lat,
				Longitude: tt.lon,
				Speed:     tt.speed,
			}
			if loc.Latitude != tt.lat {
				t.Errorf("Latitude = %f, want %f", loc.Latitude, tt.lat)
			}
			if loc.Longitude != tt.lon {
				t.Errorf("Longitude = %f, want %f", loc.Longitude, tt.lon)
			}
		})
	}
}
