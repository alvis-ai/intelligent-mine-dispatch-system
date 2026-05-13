package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aicong/mine-dispatch/pkg/utils"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type VehicleType struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name        string    `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"size:256" json:"description"`
	Icon        string    `gorm:"size:64" json:"icon"`
	Capacity    float64   `json:"capacity"`
	Weight      float64   `json:"weight"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (VehicleType) TableName() string { return "vehicle_types" }

type LoadingPoint struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Type      string    `gorm:"size:32;default:'loading'" json:"type"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Material  string    `gorm:"size:64" json:"material"`
	Status    int32     `gorm:"default:1" json:"status"`
	MineID    uint64    `gorm:"index" json:"mine_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (LoadingPoint) TableName() string { return "loading_points" }

type mgmtResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

func mgmtJSON(w http.ResponseWriter, resp mgmtResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func RegisterManagementRoutes(server *rest.Server, db *gorm.DB) {
	if db == nil {
		return
	}

	// === Vehicle Types ===
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/vehicle-types",
		Handler: listVehicleTypes(db),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/vehicle-types",
		Handler: createVehicleType(db),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/v1/vehicle-types/:id",
		Handler: updateVehicleType(db),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/v1/vehicle-types/:id",
		Handler: deleteVehicleType(db),
	})

	// === Loading Points ===
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/loading-points",
		Handler: listLoadingPoints(db),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/loading-points",
		Handler: createLoadingPoint(db),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/v1/loading-points/:id",
		Handler: updateLoadingPoint(db),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/v1/loading-points/:id",
		Handler: deleteLoadingPoint(db),
	})
}

// === Vehicle Type Handlers ===

func listVehicleTypes(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list []VehicleType
		result := db.Find(&list)
		if result.Error != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: result.Error.Error()})
			return
		}
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success", Data: list, Total: result.RowsAffected})
	}
}

func createVehicleType(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Icon        string  `json:"icon"`
			Capacity    float64 `json:"capacity"`
			Weight      float64 `json:"weight"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid request"})
			return
		}
		if req.Name == "" {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "name is required"})
			return
		}
		vt := VehicleType{
			ID:          utils.NextID(),
			Name:        req.Name,
			Description: req.Description,
			Icon:        req.Icon,
			Capacity:    req.Capacity,
			Weight:      req.Weight,
		}
		if err := db.Create(&vt).Error; err != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: err.Error()})
			return
		}
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success", Data: vt})
	}
}

func updateVehicleType(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid id"})
			return
		}
		var req struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Icon        string  `json:"icon"`
			Capacity    float64 `json:"capacity"`
			Weight      float64 `json:"weight"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid request"})
			return
		}
		updates := map[string]interface{}{}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		updates["description"] = req.Description
		updates["icon"] = req.Icon
		updates["capacity"] = req.Capacity
		updates["weight"] = req.Weight
		if err := db.Model(&VehicleType{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: err.Error()})
			return
		}
		var vt VehicleType
		db.First(&vt, id)
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success", Data: vt})
	}
}

func deleteVehicleType(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid id"})
			return
		}
		if err := db.Delete(&VehicleType{}, id).Error; err != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: err.Error()})
			return
		}
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success"})
	}
}

// === Loading Point Handlers ===

func listLoadingPoints(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list []LoadingPoint
		q := db.Model(&LoadingPoint{})
		if t := r.URL.Query().Get("type"); t != "" {
			q = q.Where("type = ?", t)
		}
		result := q.Find(&list)
		if result.Error != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: result.Error.Error()})
			return
		}
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success", Data: list, Total: result.RowsAffected})
	}
}

func createLoadingPoint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name      string  `json:"name"`
			Type      string  `json:"type"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Material  string  `json:"material"`
			Status    int32   `json:"status"`
			MineID    uint64  `json:"mine_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid request"})
			return
		}
		if req.Name == "" {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "name is required"})
			return
		}
		if req.Type == "" {
			req.Type = "loading"
		}
		lp := LoadingPoint{
			ID:        utils.NextID(),
			Name:      req.Name,
			Type:      req.Type,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Material:  req.Material,
			Status:    req.Status,
			MineID:    req.MineID,
		}
		if lp.Status == 0 {
			lp.Status = 1
		}
		if err := db.Create(&lp).Error; err != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: err.Error()})
			return
		}
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success", Data: lp})
	}
}

func updateLoadingPoint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid id"})
			return
		}
		var req struct {
			Name      string  `json:"name"`
			Type      string  `json:"type"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Material  string  `json:"material"`
			Status    int32   `json:"status"`
			MineID    uint64  `json:"mine_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid request"})
			return
		}
		updates := map[string]interface{}{}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Type != "" {
			updates["type"] = req.Type
		}
		updates["latitude"] = req.Latitude
		updates["longitude"] = req.Longitude
		updates["material"] = req.Material
		if req.Status > 0 {
			updates["status"] = req.Status
		}
		if req.MineID > 0 {
			updates["mine_id"] = req.MineID
		}
		if err := db.Model(&LoadingPoint{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: err.Error()})
			return
		}
		var lp LoadingPoint
		db.First(&lp, id)
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success", Data: lp})
	}
}

func deleteLoadingPoint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			mgmtJSON(w, mgmtResponse{Code: 400, Message: "invalid id"})
			return
		}
		if err := db.Delete(&LoadingPoint{}, id).Error; err != nil {
			mgmtJSON(w, mgmtResponse{Code: 500, Message: err.Error()})
			return
		}
		mgmtJSON(w, mgmtResponse{Code: 0, Message: "success"})
	}
}
