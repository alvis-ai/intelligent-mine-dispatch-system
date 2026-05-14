package logic

import (
	"encoding/json"
	"fmt"
	"time"

	alarmv1 "github.com/aicong/mine-dispatch/proto/alarm/v1"
	"github.com/aicong/mine-dispatch/services/alarm-service/internal/model"
	"github.com/aicong/mine-dispatch/pkg/utils"
)

// checkGeofenceViolations checks a vehicle position against all enabled geofences.
// Returns generated alarm events.
func (l *AlarmLogic) checkGeofenceViolations(lat, lon, speed float64, vehicleID uint64) []*alarmv1.AlarmEvent {
	var events []*alarmv1.AlarmEvent

	var fences []model.Geofence
	l.svc.DB.Where("enabled = ?", true).Find(&fences)
	if len(fences) == 0 {
		return nil
	}

	var rules []model.AlarmRule
	l.svc.DB.Where("enabled = ? AND rule_type = ?", true, model.RuleTypeGeofence).Find(&rules)
	if len(rules) == 0 {
		return nil
	}

	ruleMap := make(map[uint64]model.AlarmRule)
	for _, r := range rules {
		ruleMap[r.GeofenceID] = r
	}

	vehiclePlate := l.getVehiclePlate(vehicleID)

	for _, fence := range fences {
		rule, hasRule := ruleMap[fence.ID]
		if !hasRule {
			continue
		}

		var inside bool
		if fence.Shape == model.ShapeCircle {
			inside = pointInCircle(lat, lon, fence.CenterLat, fence.CenterLon, fence.RadiusM)
		} else {
			var pts []*alarmv1.Coord
			if err := json.Unmarshal([]byte(fence.PointsJSON), &pts); err == nil && len(pts) >= 3 {
				inside = pointInPolygon(lat, lon, pts)
			}
		}

		if fence.FenceType == "restricted" && inside {
			events = append(events, l.createEvent(rule, vehicleID, vehiclePlate,
				"geofence_entry", fmt.Sprintf("车辆进入禁区: %s", fence.Name), lat, lon, speed))
		}
		if fence.FenceType == "restricted" && !inside {
			// not in restricted zone — ok
		}

		// Speed check within the geofence
		if inside && fence.MaxSpeedKMH > 0 && speed > float64(fence.MaxSpeedKMH) {
			events = append(events, l.createEvent(rule, vehicleID, vehiclePlate,
				"speeding", fmt.Sprintf("在 %s 内超速: %.0f km/h (限速 %d km/h)", fence.Name, speed, fence.MaxSpeedKMH), lat, lon, speed))
		}
	}

	return events
}

// checkSpeeding checks global speeding rules (not tied to a geofence).
func (l *AlarmLogic) checkSpeeding(lat, lon, speed float64, vehicleID uint64) []*alarmv1.AlarmEvent {
	var events []*alarmv1.AlarmEvent
	if speed < 1 {
		return nil
	}

	if speed > 80 {
		var rules []model.AlarmRule
		l.svc.DB.Where("enabled = ? AND rule_type = ? AND severity = ?", true, model.RuleTypeSpeeding, model.SeverityCritical).Find(&rules)
		if len(rules) == 0 {
			return nil
		}
		vehiclePlate := l.getVehiclePlate(vehicleID)
		for _, rule := range rules {
			events = append(events, l.createEvent(rule, vehicleID, vehiclePlate,
				"speeding", fmt.Sprintf("严重超速: %.0f km/h (超过 80 km/h)", speed), lat, lon, speed))
		}
	} else if speed > 60 {
		var rules []model.AlarmRule
		l.svc.DB.Where("enabled = ? AND rule_type = ? AND severity = ?", true, model.RuleTypeSpeeding, model.SeverityWarning).Find(&rules)
		if len(rules) == 0 {
			return nil
		}
		vehiclePlate := l.getVehiclePlate(vehicleID)
		for _, rule := range rules {
			events = append(events, l.createEvent(rule, vehicleID, vehiclePlate,
				"speeding", fmt.Sprintf("超速警告: %.0f km/h", speed), lat, lon, speed))
		}
	}
	return events
}

func (l *AlarmLogic) getVehiclePlate(vehicleID uint64) string {
	type Vehicle struct {
		Plate string
	}
	var v Vehicle
	l.svc.DB.Table("vehicles").Select("plate").Where("id = ?", vehicleID).Scan(&v)
	return v.Plate
}

func (l *AlarmLogic) createEvent(rule model.AlarmRule, vehicleID uint64, plate, alarmType, message string, lat, lon, speed float64) *alarmv1.AlarmEvent {
	now := time.Now()
	event := model.AlarmEvent{
		ID:           utils.NextID(),
		RuleID:       rule.ID,
		VehicleID:    vehicleID,
		VehiclePlate: plate,
		AlarmType:    alarmType,
		Severity:     rule.Severity,
		Message:      message,
		Latitude:     lat,
		Longitude:    lon,
		Speed:        speed,
		MineID:       1,
		CreatedAt:    now,
	}
	l.svc.DB.Create(&event)

	// Publish to Redis for WebSocket push
	eventJSON, _ := json.Marshal(event)
	l.svc.Redis.Publish(l.ctx, "alarm:events", string(eventJSON))

	return &alarmv1.AlarmEvent{
		Id:           event.ID,
		RuleId:       event.RuleID,
		VehicleId:    event.VehicleID,
		VehiclePlate: event.VehiclePlate,
		AlarmType:    event.AlarmType,
		Severity:     event.Severity,
		Message:      event.Message,
		Latitude:     event.Latitude,
		Longitude:    event.Longitude,
		Speed:        event.Speed,
		Acknowledged: event.Acknowledged,
		CreatedAt:    event.CreatedAt.Format(time.RFC3339),
	}
}
