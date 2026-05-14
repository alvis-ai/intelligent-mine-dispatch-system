package logic

import (
	"math"
	"testing"

	alarmv1 "github.com/aicong/mine-dispatch/proto/alarm/v1"
)

func TestEarthDistance(t *testing.T) {
	tests := []struct {
		name string
		lat1, lon1, lat2, lon2 float64
		wantMin, wantMax float64 // acceptable range in meters
	}{
		{
			name:    "same point",
			lat1:    39.9042, lon1: 116.4074,
			lat2:    39.9042, lon2: 116.4074,
			wantMin: 0, wantMax: 1,
		},
		{
			name:    "Beijing to Shanghai (~1060 km)",
			lat1:    39.9042, lon1: 116.4074,
			lat2:    31.2304, lon2: 121.4737,
			wantMin: 1000000, wantMax: 1100000,
		},
		{
			name:    "~1 degree latitude (~111 km)",
			lat1:    40.0, lon1: 116.0,
			lat2:    41.0, lon2: 116.0,
			wantMin: 110000, wantMax: 112000,
		},
		{
			name:    "~1 degree longitude at equator (~111 km)",
			lat1:    0, lon1: 0,
			lat2:    0, lon2: 1,
			wantMin: 110000, wantMax: 112000,
		},
		{
			name:    "antipodal points (maximum)",
			lat1:    0, lon1: 0,
			lat2:    0, lon2: 180,
			wantMin: 20000000, wantMax: 20020000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := earthDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("earthDistance() = %.0f m, want [%.0f, %.0f] m", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestPointInCircle(t *testing.T) {
	centerLat, centerLon := 39.9042, 116.4074
	radiusM := 500.0

	tests := []struct {
		name string
		lat, lon float64
		want bool
	}{
		{
			name: "center point",
			lat:  centerLat, lon: centerLon,
			want: true,
		},
		{
			name: "nearby within radius (~200 m)",
			lat:  centerLat + 0.002, lon: centerLon + 0.001,
			want: true,
		},
		{
			name: "near boundary (~400 m)",
			lat:  centerLat + 0.0035, lon: centerLon,
			want: true,
		},
		{
			name: "far outside (~1 km)",
			lat:  centerLat + 0.01, lon: centerLon + 0.01,
			want: false,
		},
		{
			name: "very far",
			lat:  centerLat + 0.1, lon: centerLon + 0.1,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pointInCircle(tt.lat, tt.lon, centerLat, centerLon, radiusM)
			if got != tt.want {
				d := earthDistance(tt.lat, tt.lon, centerLat, centerLon)
				t.Errorf("pointInCircle() = %v, want %v (distance=%.0f m)", got, tt.want, d)
			}
		})
	}
}

func TestPointInPolygon(t *testing.T) {
	// A simple square polygon around (39.90, 116.40)
	polygon := []*alarmv1.Coord{
		{Lat: 39.90, Lon: 116.40},
		{Lat: 39.90, Lon: 116.41},
		{Lat: 39.91, Lon: 116.41},
		{Lat: 39.91, Lon: 116.40},
	}

	// Pentagon polygon
	pentagon := []*alarmv1.Coord{
		{Lat: 39.90, Lon: 116.40},
		{Lat: 39.905, Lon: 116.38},
		{Lat: 39.92, Lon: 116.39},
		{Lat: 39.92, Lon: 116.42},
		{Lat: 39.90, Lon: 116.42},
	}

	tests := []struct {
		name   string
		lat, lon float64
		points []*alarmv1.Coord
		want   bool
	}{
		{
			name:   "inside square",
			lat:    39.905, lon: 116.405,
			points: polygon,
			want:   true,
		},
		{
			name:   "outside square (south)",
			lat:    39.89, lon: 116.405,
			points: polygon,
			want:   false,
		},
		{
			name:   "outside square (east)",
			lat:    39.905, lon: 116.42,
			points: polygon,
			want:   false,
		},
		{
			name:   "on vertex",
			lat:    39.90, lon: 116.40,
			points: polygon,
			want:   true,
		},
		{
			name:   "inside pentagon",
			lat:    39.91, lon: 116.40,
			points: pentagon,
			want:   true,
		},
		{
			name:   "outside pentagon",
			lat:    39.89, lon: 116.40,
			points: pentagon,
			want:   false,
		},
		{
			name:   "inside pentagon concave region",
			lat:    39.908, lon: 116.39,
			points: pentagon,
			want:   true,
		},
		{
			name:   "empty polygon",
			lat:    39.905, lon: 116.405,
			points: []*alarmv1.Coord{},
			want:   false,
		},
		{
			name:   "insufficient points",
			lat:    39.905, lon: 116.405,
			points: []*alarmv1.Coord{{Lat: 39.9, Lon: 116.4}, {Lat: 39.91, Lon: 116.41}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pointInPolygon(tt.lat, tt.lon, tt.points)
			if got != tt.want {
				t.Errorf("pointInPolygon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEarthDistance_Symmetry(t *testing.T) {
	// Haversine should be symmetric
	d1 := earthDistance(39.9, 116.4, 31.2, 121.5)
	d2 := earthDistance(31.2, 121.5, 39.9, 116.4)
	if math.Abs(d1-d2) > 0.001 {
		t.Errorf("earthDistance not symmetric: %.6f vs %.6f", d1, d2)
	}
}

func TestEarthDistance_NonNegative(t *testing.T) {
	for i := 0; i < 100; i++ {
		lat1, lon1 := float64(i)*0.5, float64(i)*0.3
		lat2, lon2 := float64(i)*0.7+10, float64(i)*0.2+5
		d := earthDistance(lat1, lon1, lat2, lon2)
		if d < 0 || math.IsNaN(d) || math.IsInf(d, 0) {
			t.Errorf("earthDistance produced invalid result: %v at iteration %d", d, i)
		}
	}
}
