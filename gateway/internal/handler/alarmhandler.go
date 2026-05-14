package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	alarmv1 "github.com/aicong/mine-dispatch/proto/alarm/v1"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
)

func RegisterAlarmRoutes(server *rest.Server, conn *grpc.ClientConn) {
	client := alarmv1.NewAlarmServiceClient(conn)

	// ── Geofences ──
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/geofences",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			mineID, _ := strconv.ParseUint(r.URL.Query().Get("mine_id"), 10, 64)
			fenceType := r.URL.Query().Get("fence_type")
			enabledOnly := r.URL.Query().Get("enabled_only") == "true"
			resp, _ := client.ListGeofences(r.Context(), &alarmv1.ListGeofencesRequest{
				MineId:      mineID,
				FenceType:   fenceType,
				EnabledOnly: enabledOnly,
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/geofences",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req alarmv1.CreateGeofenceRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CreateGeofence(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/geofences/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.GetGeofence(r.Context(), &alarmv1.GetGeofenceRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPut,
		Path:   "/api/v1/geofences/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			body, _ := io.ReadAll(r.Body)
			var req alarmv1.UpdateGeofenceRequest
			json.Unmarshal(body, &req)
			req.Id = id
			resp, _ := client.UpdateGeofence(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodDelete,
		Path:   "/api/v1/geofences/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.DeleteGeofence(r.Context(), &alarmv1.DeleteGeofenceRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	// ── Alarm Rules ──
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/alarm-rules",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			rtype := r.URL.Query().Get("rule_type")
			enabledOnly := r.URL.Query().Get("enabled_only") == "true"
			resp, _ := client.ListAlarmRules(r.Context(), &alarmv1.ListAlarmRulesRequest{
				RuleType:    rtype,
				EnabledOnly: enabledOnly,
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/alarm-rules",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req alarmv1.CreateAlarmRuleRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CreateAlarmRule(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPut,
		Path:   "/api/v1/alarm-rules/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			body, _ := io.ReadAll(r.Body)
			var req alarmv1.UpdateAlarmRuleRequest
			json.Unmarshal(body, &req)
			req.Id = id
			resp, _ := client.UpdateAlarmRule(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodDelete,
		Path:   "/api/v1/alarm-rules/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			resp, _ := client.DeleteAlarmRule(r.Context(), &alarmv1.DeleteAlarmRuleRequest{Id: id})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	// ── Alarm Events ──
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/alarms",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			vehicleID, _ := strconv.ParseUint(q.Get("vehicle_id"), 10, 64)
			page, _ := strconv.Atoi(q.Get("page"))
			pageSize, _ := strconv.Atoi(q.Get("page_size"))
			mineID, _ := strconv.ParseUint(q.Get("mine_id"), 10, 64)
			if page < 1 {
				page = 1
			}
			if pageSize < 1 {
				pageSize = 20
			}
			resp, _ := client.ListAlarmEvents(r.Context(), &alarmv1.ListAlarmEventsRequest{
				VehicleId:          vehicleID,
				Severity:           q.Get("severity"),
				AlarmType:          q.Get("alarm_type"),
				UnacknowledgedOnly: q.Get("unacknowledged_only") == "true",
				MineId:             mineID,
				Page:               int32(page),
				PageSize:           int32(pageSize),
			})
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/alarms/:id/acknowledge",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
			body, _ := io.ReadAll(r.Body)
			var req alarmv1.AcknowledgeAlarmRequest
			json.Unmarshal(body, &req)
			req.Id = id
			resp, _ := client.AcknowledgeAlarm(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	// ── Dashboard stats (aggregated from alarms) ──
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/alarms/stats",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// Unacknowledged critical
			criticalResp, _ := client.ListAlarmEvents(r.Context(), &alarmv1.ListAlarmEventsRequest{
				Severity:           "critical",
				UnacknowledgedOnly: true,
				PageSize:           1,
			})
			warningResp, _ := client.ListAlarmEvents(r.Context(), &alarmv1.ListAlarmEventsRequest{
				Severity:           "warning",
				UnacknowledgedOnly: true,
				PageSize:           1,
			})

			stats := map[string]interface{}{
				"unacknowledged_critical": criticalResp.Total,
				"unacknowledged_warning":  warningResp.Total,
				"total_critical":          0,
				"total_warning":           0,
			}
			allCritical, _ := client.ListAlarmEvents(r.Context(), &alarmv1.ListAlarmEventsRequest{
				Severity: "critical",
				PageSize: 1,
			})
			allWarning, _ := client.ListAlarmEvents(r.Context(), &alarmv1.ListAlarmEventsRequest{
				Severity: "warning",
				PageSize: 1,
			})
			stats["total_critical"] = allCritical.Total
			stats["total_warning"] = allWarning.Total

			resp := map[string]interface{}{"code": 0, "message": "success", "data": stats}
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	// WebSocket alarm subscription via existing WS hub
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/alarms/recent",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			resp, _ := client.ListAlarmEvents(r.Context(), &alarmv1.ListAlarmEventsRequest{
				UnacknowledgedOnly: true,
				Page:              1,
				PageSize:          10,
			})
			// If no unack, return latest 10
			if resp.Total == 0 {
				resp, _ = client.ListAlarmEvents(r.Context(), &alarmv1.ListAlarmEventsRequest{
					Page:     1,
					PageSize: 10,
				})
			}
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	// ── Position check (telemetry integration) ──
	server.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/api/v1/alarms/check-position",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req alarmv1.CheckPositionRequest
			json.Unmarshal(body, &req)
			resp, _ := client.CheckPosition(r.Context(), &req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})
}

