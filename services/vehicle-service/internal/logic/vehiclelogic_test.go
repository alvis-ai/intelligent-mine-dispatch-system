package logic

import (
	"testing"

	vehiclev1 "github.com/aicong/mine-dispatch/proto/vehicle/v1"
)

func TestVehicleTypeValues(t *testing.T) {
	tests := []struct {
		name     string
		vt       vehiclev1.VehicleType
		wantName string
	}{
		{"unspecified", vehiclev1.VehicleType_VEHICLE_TYPE_UNSPECIFIED, "VEHICLE_TYPE_UNSPECIFIED"},
		{"mining truck", vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK, "VEHICLE_TYPE_MINING_TRUCK"},
		{"excavator", vehiclev1.VehicleType_VEHICLE_TYPE_EXCAVATOR, "VEHICLE_TYPE_EXCAVATOR"},
		{"loader", vehiclev1.VehicleType_VEHICLE_TYPE_LOADER, "VEHICLE_TYPE_LOADER"},
		{"bulldozer", vehiclev1.VehicleType_VEHICLE_TYPE_BULLDOZER, "VEHICLE_TYPE_BULLDOZER"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := vehiclev1.VehicleType_name[int32(tt.vt)]
			if name != tt.wantName {
				t.Errorf("VehicleType_name[%d] = %s, want %s", tt.vt, name, tt.wantName)
			}
		})
	}
}

func TestVehicleStatusValues(t *testing.T) {
	tests := []struct {
		name     string
		vs       vehiclev1.VehicleStatus
		wantName string
	}{
		{"unspecified", vehiclev1.VehicleStatus_VEHICLE_STATUS_UNSPECIFIED, "VEHICLE_STATUS_UNSPECIFIED"},
		{"idle", vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE, "VEHICLE_STATUS_IDLE"},
		{"loading", vehiclev1.VehicleStatus_VEHICLE_STATUS_LOADING, "VEHICLE_STATUS_LOADING"},
		{"maintenance", vehiclev1.VehicleStatus_VEHICLE_STATUS_MAINTENANCE, "VEHICLE_STATUS_MAINTENANCE"},
		{"offline", vehiclev1.VehicleStatus_VEHICLE_STATUS_OFFLINE, "VEHICLE_STATUS_OFFLINE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := vehiclev1.VehicleStatus_name[int32(tt.vs)]
			if name != tt.wantName {
				t.Errorf("VehicleStatus_name[%d] = %s, want %s", tt.vs, name, tt.wantName)
			}
		})
	}
}

func TestVehicleResponse_Success(t *testing.T) {
	resp := &vehiclev1.VehicleResponse{
		Code:    0,
		Message: "success",
		Data: &vehiclev1.Vehicle{
			Id:        1001,
			Plate:     "矿卡-A001",
			Type:      vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK,
			Status:    vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE,
			Latitude:  39.9042,
			Longitude: 116.4074,
			FuelLevel: 85.5,
			MineId:    1,
		},
	}
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if resp.Data.Plate != "矿卡-A001" {
		t.Errorf("Plate = %s, want 矿卡-A001", resp.Data.Plate)
	}
	if resp.Data.Type != vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK {
		t.Errorf("Type = %v, want MINING_TRUCK", resp.Data.Type)
	}
	if resp.Data.Status != vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE {
		t.Errorf("Status = %v, want IDLE", resp.Data.Status)
	}
	if resp.Data.FuelLevel != 85.5 {
		t.Errorf("FuelLevel = %f, want 85.5", resp.Data.FuelLevel)
	}
}

func TestVehicleResponse_NotFound(t *testing.T) {
	resp := &vehiclev1.VehicleResponse{Code: 404, Message: "vehicle not found"}
	if resp.Code != 404 {
		t.Errorf("Code = %d, want 404", resp.Code)
	}
	if resp.Data != nil {
		t.Error("Data should be nil when not found")
	}
}

func TestVehicleListResponse_Pagination(t *testing.T) {
	resp := &vehiclev1.VehicleListResponse{
		Code:    0,
		Message: "success",
		Data: []*vehiclev1.Vehicle{
			{Id: 1, Plate: "矿卡-A001", Type: vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK, Status: vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE},
			{Id: 2, Plate: "挖机-B001", Type: vehiclev1.VehicleType_VEHICLE_TYPE_EXCAVATOR, Status: vehiclev1.VehicleStatus_VEHICLE_STATUS_LOADING},
		},
		Total: 10,
	}
	if resp.Total != 10 {
		t.Errorf("Total = %d, want 10", resp.Total)
	}
	if len(resp.Data) != 2 {
		t.Errorf("len(Data) = %d, want 2", len(resp.Data))
	}
	if resp.Data[0].Type != vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK {
		t.Errorf("Data[0].Type = %v, want MINING_TRUCK", resp.Data[0].Type)
	}
	if resp.Data[1].Status != vehiclev1.VehicleStatus_VEHICLE_STATUS_LOADING {
		t.Errorf("Data[1].Status = %v, want LOADING", resp.Data[1].Status)
	}
}

func TestVehicleListResponse_Empty(t *testing.T) {
	resp := &vehiclev1.VehicleListResponse{
		Code:    0,
		Message: "success",
		Data:    []*vehiclev1.Vehicle{},
		Total:   0,
	}
	if len(resp.Data) != 0 {
		t.Errorf("len(Data) = %d, want 0", len(resp.Data))
	}
}

func TestCreateVehicleRequest_DefaultStatus(t *testing.T) {
	req := &vehiclev1.CreateVehicleRequest{
		Plate:  "矿卡-A003",
		Type:   vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK,
		MineId: 1,
	}
	if req.Plate == "" {
		t.Error("Plate should not be empty")
	}
	if req.MineId == 0 {
		t.Error("MineId should not be zero")
	}
}

func TestVehicleTypeEnum_String(t *testing.T) {
	if vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK.String() != "VEHICLE_TYPE_MINING_TRUCK" {
		t.Errorf("Unexpected string: %s", vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK.String())
	}
}

func TestVehicle_AllTypes(t *testing.T) {
	types := []vehiclev1.VehicleType{
		vehiclev1.VehicleType_VEHICLE_TYPE_UNSPECIFIED,
		vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK,
		vehiclev1.VehicleType_VEHICLE_TYPE_EXCAVATOR,
		vehiclev1.VehicleType_VEHICLE_TYPE_LOADER,
		vehiclev1.VehicleType_VEHICLE_TYPE_BULLDOZER,
	}

	for _, vt := range types {
		t.Run(vt.String(), func(t *testing.T) {
			v := &vehiclev1.Vehicle{
				Id:     1,
				Plate:  "test",
				Type:   vt,
				Status: vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE,
			}
			if v.Type != vt {
				t.Errorf("Type = %v, want %v", v.Type, vt)
			}
		})
	}
}

func TestVehicle_AllStatuses(t *testing.T) {
	statuses := []vehiclev1.VehicleStatus{
		vehiclev1.VehicleStatus_VEHICLE_STATUS_UNSPECIFIED,
		vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE,
		vehiclev1.VehicleStatus_VEHICLE_STATUS_LOADING,
		vehiclev1.VehicleStatus_VEHICLE_STATUS_MAINTENANCE,
		vehiclev1.VehicleStatus_VEHICLE_STATUS_OFFLINE,
	}

	for _, vs := range statuses {
		t.Run(vs.String(), func(t *testing.T) {
			v := &vehiclev1.Vehicle{
				Id:     1,
				Plate:  "test",
				Type:   vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK,
				Status: vs,
			}
			if v.Status != vs {
				t.Errorf("Status = %v, want %v", v.Status, vs)
			}
		})
	}
}

func TestListVehicleRequest_Filters(t *testing.T) {
	tests := []struct {
		name   string
		req    *vehiclev1.ListVehicleRequest
		hasF   bool
	}{
		{"no filter", &vehiclev1.ListVehicleRequest{Page: 1, PageSize: 10}, false},
		{"by type", &vehiclev1.ListVehicleRequest{Page: 1, PageSize: 10, Type: vehiclev1.VehicleType_VEHICLE_TYPE_MINING_TRUCK}, true},
		{"by status", &vehiclev1.ListVehicleRequest{Page: 1, PageSize: 10, Status: vehiclev1.VehicleStatus_VEHICLE_STATUS_IDLE}, true},
		{"by mine", &vehiclev1.ListVehicleRequest{Page: 1, PageSize: 10, MineId: 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasFilter := tt.req.Type != vehiclev1.VehicleType_VEHICLE_TYPE_UNSPECIFIED ||
				tt.req.Status != vehiclev1.VehicleStatus_VEHICLE_STATUS_UNSPECIFIED ||
				tt.req.MineId > 0
			if hasFilter != tt.hasF {
				t.Errorf("hasFilter = %v, want %v", hasFilter, tt.hasF)
			}
		})
	}
}
