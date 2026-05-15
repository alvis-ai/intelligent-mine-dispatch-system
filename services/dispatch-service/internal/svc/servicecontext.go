package svc

import (
	"github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/config"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	RouteClient routev1.RouteServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.PostgresDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPass,
	})

	var routeClient routev1.RouteServiceClient
	if c.RouteSvcAddr != "" {
		conn, err := grpc.Dial(c.RouteSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		routeClient = routev1.NewRouteServiceClient(conn)
	}

	return &ServiceContext{Config: c, DB: db, Redis: rdb, RouteClient: routeClient}
}
