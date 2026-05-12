package main

import (
	"flag"
	"fmt"

	"github.com/aicong/mine-dispatch/proto/vehicle/v1"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/config"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/server"
	"github.com/aicong/mine-dispatch/services/vehicle-service/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/vehicle.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.SetLevel(logx.InfoLevel)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		vehiclev1.RegisterVehicleServiceServer(grpcServer, server.NewVehicleServer(ctx))
		if c.RpcServerConf.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	fmt.Printf("Starting vehicle-service on %s\n", c.RpcServerConf.ListenOn)
	s.Start()
}
