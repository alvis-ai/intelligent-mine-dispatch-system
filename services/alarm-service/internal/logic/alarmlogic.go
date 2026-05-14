package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	alarmv1 "github.com/aicong/mine-dispatch/proto/alarm/v1"
	"github.com/aicong/mine-dispatch/services/alarm-service/internal/model"
	"github.com/aicong/mine-dispatch/services/alarm-service/internal/svc"
	"github.com/aicong/mine-dispatch/pkg/utils"
)

type AlarmLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewAlarmLogic(ctx context.Context, svc *svc.ServiceContext) *AlarmLogic {
	return &AlarmLogic{ctx: ctx, svc: svc}
}

// ── Geofence CRUD ──

func (l *AlarmLogic) CreateGeofence(in *alarmv1.CreateGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	ptsJSON, _ := json.Marshal(in.Points)
	g := model.Geofence{
		ID:          utils.NextID(),
		Name:        in.Name,
		Shape:       in.Shape,
		CenterLat:   in.CenterLat,
		CenterLon:   in.CenterLon,
		RadiusM:     in.RadiusM,
		PointsJSON:  string(ptsJSON),
		FenceType:   in.FenceType,
		MinSpeedKMH: in.MinSpeedKmh,
		MaxSpeedKMH: in.MaxSpeedKmh,
		TimeRange:   in.TimeRange,
		Enabled:     in.Enabled,
		MineID:      in.MineId,
	}
	if err := l.svc.DB.Create(&g).Error; err != nil {
		return &alarmv1.GeofenceResponse{Code: 500, Message: err.Error()}, nil
	}
	return &alarmv1.GeofenceResponse{Code: 0, Message: "success", Data: modelToProtoGeofence(&g)}, nil
}

func (l *AlarmLogic) UpdateGeofence(in *alarmv1.UpdateGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	var g model.Geofence
	if err := l.svc.DB.First(&g, in.Id).Error; err != nil {
		return &alarmv1.GeofenceResponse{Code: 404, Message: "geofence not found"}, nil
	}
	updates := map[string]interface{}{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Shape != "" {
		updates["shape"] = in.Shape
	}
	if in.CenterLat != 0 || in.CenterLon != 0 {
		updates["center_lat"] = in.CenterLat
		updates["center_lon"] = in.CenterLon
	}
	if in.RadiusM != 0 {
		updates["radius_m"] = in.RadiusM
	}
	if in.Points != nil {
		ptsJSON, _ := json.Marshal(in.Points)
		updates["points_json"] = string(ptsJSON)
	}
	if in.FenceType != "" {
		updates["fence_type"] = in.FenceType
	}
	if in.MinSpeedKmh != 0 {
		updates["min_speed_kmh"] = in.MinSpeedKmh
	}
	if in.MaxSpeedKmh != 0 {
		updates["max_speed_kmh"] = in.MaxSpeedKmh
	}
	updates["enabled"] = in.Enabled
	if in.MineId != 0 {
		updates["mine_id"] = in.MineId
	}
	l.svc.DB.Model(&g).Updates(updates)
	l.svc.DB.First(&g, in.Id)
	return &alarmv1.GeofenceResponse{Code: 0, Message: "success", Data: modelToProtoGeofence(&g)}, nil
}

func (l *AlarmLogic) GetGeofence(in *alarmv1.GetGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	var g model.Geofence
	if err := l.svc.DB.First(&g, in.Id).Error; err != nil {
		return &alarmv1.GeofenceResponse{Code: 404, Message: "geofence not found"}, nil
	}
	return &alarmv1.GeofenceResponse{Code: 0, Message: "success", Data: modelToProtoGeofence(&g)}, nil
}

func (l *AlarmLogic) DeleteGeofence(in *alarmv1.DeleteGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	l.svc.DB.Delete(&model.Geofence{}, in.Id)
	l.svc.DB.Where("geofence_id = ?", in.Id).Delete(&model.AlarmRule{})
	return &alarmv1.GeofenceResponse{Code: 0, Message: "deleted"}, nil
}

func (l *AlarmLogic) ListGeofences(in *alarmv1.ListGeofencesRequest) (*alarmv1.GeofenceListResponse, error) {
	var fences []model.Geofence
	var total int64
	db := l.svc.DB.Model(&model.Geofence{})
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	if in.FenceType != "" {
		db = db.Where("fence_type = ?", in.FenceType)
	}
	if in.EnabledOnly {
		db = db.Where("enabled = ?", true)
	}
	db.Count(&total)
	db.Order("id ASC").Find(&fences)
	var list []*alarmv1.Geofence
	for i := range fences {
		list = append(list, modelToProtoGeofence(&fences[i]))
	}
	return &alarmv1.GeofenceListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}

// ── Alarm Rule CRUD ──

func (l *AlarmLogic) CreateAlarmRule(in *alarmv1.CreateAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	r := model.AlarmRule{
		ID:          utils.NextID(),
		Name:        in.Name,
		RuleType:    in.RuleType,
		GeofenceID:  in.GeofenceId,
		Severity:    in.Severity,
		Description: in.Description,
		Enabled:     in.Enabled,
	}
	if err := l.svc.DB.Create(&r).Error; err != nil {
		return &alarmv1.AlarmRuleResponse{Code: 500, Message: err.Error()}, nil
	}
	return &alarmv1.AlarmRuleResponse{Code: 0, Message: "success", Data: modelToProtoRule(&r)}, nil
}

func (l *AlarmLogic) UpdateAlarmRule(in *alarmv1.UpdateAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	var r model.AlarmRule
	if err := l.svc.DB.First(&r, in.Id).Error; err != nil {
		return &alarmv1.AlarmRuleResponse{Code: 404, Message: "rule not found"}, nil
	}
	updates := map[string]interface{}{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.RuleType != "" {
		updates["rule_type"] = in.RuleType
	}
	if in.GeofenceId != 0 {
		updates["geofence_id"] = in.GeofenceId
	}
	if in.Severity != "" {
		updates["severity"] = in.Severity
	}
	if in.Description != "" {
		updates["description"] = in.Description
	}
	updates["enabled"] = in.Enabled
	l.svc.DB.Model(&r).Updates(updates)
	l.svc.DB.First(&r, in.Id)
	return &alarmv1.AlarmRuleResponse{Code: 0, Message: "success", Data: modelToProtoRule(&r)}, nil
}

func (l *AlarmLogic) GetAlarmRule(in *alarmv1.GetAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	var r model.AlarmRule
	if err := l.svc.DB.First(&r, in.Id).Error; err != nil {
		return &alarmv1.AlarmRuleResponse{Code: 404, Message: "rule not found"}, nil
	}
	return &alarmv1.AlarmRuleResponse{Code: 0, Message: "success", Data: modelToProtoRule(&r)}, nil
}

func (l *AlarmLogic) DeleteAlarmRule(in *alarmv1.DeleteAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	l.svc.DB.Delete(&model.AlarmRule{}, in.Id)
	return &alarmv1.AlarmRuleResponse{Code: 0, Message: "deleted"}, nil
}

func (l *AlarmLogic) ListAlarmRules(in *alarmv1.ListAlarmRulesRequest) (*alarmv1.AlarmRuleListResponse, error) {
	var rules []model.AlarmRule
	var total int64
	db := l.svc.DB.Model(&model.AlarmRule{})
	if in.RuleType != "" {
		db = db.Where("rule_type = ?", in.RuleType)
	}
	if in.EnabledOnly {
		db = db.Where("enabled = ?", true)
	}
	db.Count(&total)
	db.Order("id ASC").Find(&rules)
	var list []*alarmv1.AlarmRule
	for i := range rules {
		list = append(list, modelToProtoRule(&rules[i]))
	}
	return &alarmv1.AlarmRuleListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}

// ── Alarm Events ──

func (l *AlarmLogic) ListAlarmEvents(in *alarmv1.ListAlarmEventsRequest) (*alarmv1.AlarmEventListResponse, error) {
	var events []model.AlarmEvent
	var total int64
	db := l.svc.DB.Model(&model.AlarmEvent{})
	if in.VehicleId > 0 {
		db = db.Where("vehicle_id = ?", in.VehicleId)
	}
	if in.Severity != "" {
		db = db.Where("severity = ?", in.Severity)
	}
	if in.AlarmType != "" {
		db = db.Where("alarm_type = ?", in.AlarmType)
	}
	if in.UnacknowledgedOnly {
		db = db.Where("acknowledged = ?", false)
	}
	if in.MineId > 0 {
		db = db.Where("mine_id = ?", in.MineId)
	}
	page := int(in.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	db.Count(&total)
	db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&events)
	var list []*alarmv1.AlarmEvent
	for _, e := range events {
		list = append(list, modelToProtoEvent(&e))
	}
	return &alarmv1.AlarmEventListResponse{Code: 0, Message: "success", Data: list, Total: total}, nil
}

func (l *AlarmLogic) AcknowledgeAlarm(in *alarmv1.AcknowledgeAlarmRequest) (*alarmv1.AlarmEventResponse, error) {
	var e model.AlarmEvent
	if err := l.svc.DB.First(&e, in.Id).Error; err != nil {
		return &alarmv1.AlarmEventResponse{Code: 404, Message: "event not found"}, nil
	}
	now := time.Now()
	l.svc.DB.Model(&e).Updates(map[string]interface{}{
		"acknowledged":    true,
		"acknowledged_by": in.AcknowledgedBy,
		"acknowledged_at": now,
	})
	e.Acknowledged = true
	e.AcknowledgedBy = in.AcknowledgedBy
	e.AcknowledgedAt = &now
	return &alarmv1.AlarmEventResponse{Code: 0, Message: "success", Data: modelToProtoEvent(&e)}, nil
}

// ── Position Check (called by telemetry on GPS report) ──

func (l *AlarmLogic) CheckPosition(in *alarmv1.CheckPositionRequest) (*alarmv1.CheckPositionResponse, error) {
	var events []*alarmv1.AlarmEvent
	geoEvents := l.checkGeofenceViolations(in.Latitude, in.Longitude, in.Speed, in.VehicleId)
	events = append(events, geoEvents...)
	speedEvents := l.checkSpeeding(in.Latitude, in.Longitude, in.Speed, in.VehicleId)
	events = append(events, speedEvents...)
	if len(events) > 0 {
		fmt.Printf("Alarm generated: %d events for vehicle %d\n", len(events), in.VehicleId)
	}
	return &alarmv1.CheckPositionResponse{Code: 0, Message: "success", Alarms: events}, nil
}

// ── Model <-> Proto conversions ──

func modelToProtoGeofence(g *model.Geofence) *alarmv1.Geofence {
	var pts []*alarmv1.Coord
	json.Unmarshal([]byte(g.PointsJSON), &pts)
	return &alarmv1.Geofence{
		Id:          g.ID,
		Name:        g.Name,
		Shape:       g.Shape,
		CenterLat:   g.CenterLat,
		CenterLon:   g.CenterLon,
		RadiusM:     g.RadiusM,
		Points:      pts,
		FenceType:   g.FenceType,
		MinSpeedKmh: g.MinSpeedKMH,
		MaxSpeedKmh: g.MaxSpeedKMH,
		TimeRange:   g.TimeRange,
		Enabled:     g.Enabled,
		MineId:      g.MineID,
		CreatedAt:   g.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   g.UpdatedAt.Format(time.RFC3339),
	}
}

func modelToProtoRule(r *model.AlarmRule) *alarmv1.AlarmRule {
	return &alarmv1.AlarmRule{
		Id:          r.ID,
		Name:        r.Name,
		RuleType:    r.RuleType,
		GeofenceId:  r.GeofenceID,
		Severity:    r.Severity,
		Description: r.Description,
		Enabled:     r.Enabled,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
	}
}

func modelToProtoEvent(e *model.AlarmEvent) *alarmv1.AlarmEvent {
	ackAt := ""
	if e.AcknowledgedAt != nil {
		ackAt = e.AcknowledgedAt.Format(time.RFC3339)
	}
	return &alarmv1.AlarmEvent{
		Id:              e.ID,
		RuleId:          e.RuleID,
		VehicleId:       e.VehicleID,
		VehiclePlate:    e.VehiclePlate,
		AlarmType:       e.AlarmType,
		Severity:        e.Severity,
		Message:         e.Message,
		Latitude:        e.Latitude,
		Longitude:       e.Longitude,
		Speed:           e.Speed,
		Acknowledged:    e.Acknowledged,
		AcknowledgedBy:  e.AcknowledgedBy,
		AcknowledgedAt:  ackAt,
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
		MineId:          e.MineID,
	}
}
