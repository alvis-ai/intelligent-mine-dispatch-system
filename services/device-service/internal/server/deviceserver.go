package server

import (
	"context"

	devicev1 "github.com/aicong/mine-dispatch/proto/device/v1"
	"github.com/aicong/mine-dispatch/services/device-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/device-service/internal/svc"
)

type DeviceServer struct {
	svc *svc.ServiceContext
	devicev1.UnimplementedDeviceServiceServer
}

func NewDeviceServer(svc *svc.ServiceContext) *DeviceServer {
	return &DeviceServer{svc: svc}
}

func (s *DeviceServer) CreateDevice(ctx context.Context, in *devicev1.CreateDeviceRequest) (*devicev1.DeviceResponse, error) {
	return logic.NewCreateDeviceLogic(ctx, s.svc).CreateDevice(in)
}

func (s *DeviceServer) GetDevice(ctx context.Context, in *devicev1.GetDeviceRequest) (*devicev1.DeviceResponse, error) {
	return logic.NewGetDeviceLogic(ctx, s.svc).GetDevice(in)
}

func (s *DeviceServer) UpdateDevice(ctx context.Context, in *devicev1.UpdateDeviceRequest) (*devicev1.DeviceResponse, error) {
	return logic.NewUpdateDeviceLogic(ctx, s.svc).UpdateDevice(in)
}

func (s *DeviceServer) DeleteDevice(ctx context.Context, in *devicev1.DeleteDeviceRequest) (*devicev1.DeviceResponse, error) {
	return logic.NewDeleteDeviceLogic(ctx, s.svc).DeleteDevice(in)
}

func (s *DeviceServer) ListDevice(ctx context.Context, in *devicev1.ListDeviceRequest) (*devicev1.DeviceListResponse, error) {
	return logic.NewListDeviceLogic(ctx, s.svc).ListDevice(in)
}
