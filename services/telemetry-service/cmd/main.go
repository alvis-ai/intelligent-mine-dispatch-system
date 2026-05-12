package main

import (
	"flag"
	"fmt"

	"github.com/aicong/mine-dispatch/proto/telemetry/v1"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/config"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/server"
	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/telemetry.yaml", "config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.SetLevel(logx.InfoLevel)
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		telemetryv1.RegisterTelemetryServiceServer(grpcServer, server.NewTelemetryServer(ctx))
		if c.RpcServerConf.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	fmt.Printf("Starting telemetry-service on %s\n", c.RpcServerConf.ListenOn)
	s.Start()
}
