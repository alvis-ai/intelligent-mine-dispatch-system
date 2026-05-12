package server

import (
	"context"

	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/svc"
)

type VehicleServer struct {
	svc *svc.ServiceContext
	vehiclev1.UnimplementedVehicleServiceServer
}

func NewVehicleServer(svc *svc.ServiceContext) *VehicleServer {
	return &VehicleServer{svc: svc}
}

func (s *VehicleServer) CreateVehicle(ctx context.Context, in *vehiclev1.CreateVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	return logic.NewCreateVehicleLogic(ctx, s.svc).CreateVehicle(in)
}
func (s *VehicleServer) GetVehicle(ctx context.Context, in *vehiclev1.GetVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	return logic.NewGetVehicleLogic(ctx, s.svc).GetVehicle(in)
}
func (s *VehicleServer) UpdateVehicle(ctx context.Context, in *vehiclev1.UpdateVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	return logic.NewUpdateVehicleLogic(ctx, s.svc).UpdateVehicle(in)
}
func (s *VehicleServer) DeleteVehicle(ctx context.Context, in *vehiclev1.DeleteVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	return logic.NewDeleteVehicleLogic(ctx, s.svc).DeleteVehicle(in)
}
func (s *VehicleServer) ListVehicle(ctx context.Context, in *vehiclev1.ListVehicleRequest) (*vehiclev1.VehicleListResponse, error) {
	return logic.NewListVehicleLogic(ctx, s.svc).ListVehicle(in)
}
