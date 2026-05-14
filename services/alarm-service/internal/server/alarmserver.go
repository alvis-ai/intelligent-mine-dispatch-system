package server

import (
	"context"

	alarmv1 "github.com/aicong/mine-dispatch/proto/alarm/v1"
	"github.com/aicong/mine-dispatch/services/alarm-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/alarm-service/internal/svc"
)

type AlarmServer struct {
	svc *svc.ServiceContext
	alarmv1.UnimplementedAlarmServiceServer
}

func NewAlarmServer(svc *svc.ServiceContext) *AlarmServer {
	return &AlarmServer{svc: svc}
}

func (s *AlarmServer) CreateGeofence(ctx context.Context, in *alarmv1.CreateGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).CreateGeofence(in)
}
func (s *AlarmServer) UpdateGeofence(ctx context.Context, in *alarmv1.UpdateGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).UpdateGeofence(in)
}
func (s *AlarmServer) GetGeofence(ctx context.Context, in *alarmv1.GetGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).GetGeofence(in)
}
func (s *AlarmServer) DeleteGeofence(ctx context.Context, in *alarmv1.DeleteGeofenceRequest) (*alarmv1.GeofenceResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).DeleteGeofence(in)
}
func (s *AlarmServer) ListGeofences(ctx context.Context, in *alarmv1.ListGeofencesRequest) (*alarmv1.GeofenceListResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).ListGeofences(in)
}

func (s *AlarmServer) CreateAlarmRule(ctx context.Context, in *alarmv1.CreateAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).CreateAlarmRule(in)
}
func (s *AlarmServer) UpdateAlarmRule(ctx context.Context, in *alarmv1.UpdateAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).UpdateAlarmRule(in)
}
func (s *AlarmServer) GetAlarmRule(ctx context.Context, in *alarmv1.GetAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).GetAlarmRule(in)
}
func (s *AlarmServer) DeleteAlarmRule(ctx context.Context, in *alarmv1.DeleteAlarmRuleRequest) (*alarmv1.AlarmRuleResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).DeleteAlarmRule(in)
}
func (s *AlarmServer) ListAlarmRules(ctx context.Context, in *alarmv1.ListAlarmRulesRequest) (*alarmv1.AlarmRuleListResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).ListAlarmRules(in)
}

func (s *AlarmServer) ListAlarmEvents(ctx context.Context, in *alarmv1.ListAlarmEventsRequest) (*alarmv1.AlarmEventListResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).ListAlarmEvents(in)
}
func (s *AlarmServer) AcknowledgeAlarm(ctx context.Context, in *alarmv1.AcknowledgeAlarmRequest) (*alarmv1.AlarmEventResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).AcknowledgeAlarm(in)
}
func (s *AlarmServer) CheckPosition(ctx context.Context, in *alarmv1.CheckPositionRequest) (*alarmv1.CheckPositionResponse, error) {
	return logic.NewAlarmLogic(ctx, s.svc).CheckPosition(in)
}
