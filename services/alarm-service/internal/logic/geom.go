package logic

import (
	"math"

	alarmv1 "github.com/aicong/mine-dispatch/proto/alarm/v1"
)

// earthDistance computes the Haversine distance in meters between two lat/lon points.
func earthDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// pointInCircle checks if a point is within a circular geofence.
func pointInCircle(lat, lon, centerLat, centerLon, radiusM float64) bool {
	return earthDistance(lat, lon, centerLat, centerLon) <= radiusM
}

// pointInPolygon uses the ray-casting algorithm to check if a point is inside a polygon.
func pointInPolygon(lat, lon float64, points []*alarmv1.Coord) bool {
	n := len(points)
	if n < 3 {
		return false
	}
	inside := false
	j := n - 1
	for i := 0; i < n; i++ {
		pi, pj := points[i], points[j]
		if (pi.Lon > lon) != (pj.Lon > lon) &&
			lat < (pj.Lat-pi.Lat)*(lon-pi.Lon)/(pj.Lon-pi.Lon)+pi.Lat {
			inside = !inside
		}
		j = i
	}
	return inside
}
