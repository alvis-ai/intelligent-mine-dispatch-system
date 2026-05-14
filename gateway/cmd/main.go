package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/aicong/mine-dispatch/gateway/internal/handler"
	userv1 "github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

type GatewayConfig struct {
	rest.RestConf
	RedisAddr    string
	RedisPass    string
	AuthSvcAddr  string
	UserSvcAddr      string
	VehicleSvcAddr   string
	TelemetrySvcAddr string
	DispatchSvcAddr  string
	AlarmSvcAddr     string
	RouteSvcAddr     string
	PostgresDSN      string
}

func proxyHandler(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()

		req, _ := http.NewRequest(r.Method, targetURL, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 502, "message": "auth service unavailable: " + err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
	}
}

func registerAuthProxy(server *rest.Server, authAddr string) {
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/auth/login",
		Handler: proxyHandler("http://" + authAddr + "/api/v1/auth/login"),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/auth/validate",
		Handler: proxyHandler("http://" + authAddr + "/api/v1/auth/validate"),
	})
}

func main() {
	flag.Parse()
	var c GatewayConfig
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	userConn, _ := grpc.Dial(c.UserSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	vehicleConn, _ := grpc.Dial(c.VehicleSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	telemetryConn, _ := grpc.Dial(c.TelemetrySvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	dispatchConn, _ := grpc.Dial(c.DispatchSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	alarmConn, _ := grpc.Dial(c.AlarmSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	routeConn, _ := grpc.Dial(c.RouteSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	userClient := userv1.NewUserServiceClient(userConn)

	rdb := redis.NewClient(&redis.Options{
		Addr: c.RedisAddr,
		Password: c.RedisPass,
	})

	handler.RegisterAuthRoutes(server, userClient)
	registerAuthProxy(server, c.AuthSvcAddr)
	handler.RegisterVehicleRoutes(server, vehicleConn)
	handler.RegisterTelemetryRoutes(server, telemetryConn)
	handler.RegisterDispatchRoutes(server, dispatchConn)
	wsHub := handler.NewWSHub(rdb)
	handler.RegisterWSRoute(server, wsHub)

	// Init GORM for management APIs
	var mgmtDB *gorm.DB
	if c.PostgresDSN != "" {
		var err error
		mgmtDB, err = gorm.Open(postgres.Open(c.PostgresDSN), &gorm.Config{})
		if err != nil {
			fmt.Printf("Warning: failed to connect management DB: %v\n", err)
		}
	}
	handler.RegisterAlarmRoutes(server, alarmConn)
	handler.RegisterRouteRoutes(server, routeConn)
	handler.RegisterManagementRoutes(server, mgmtDB)

	fmt.Printf("Starting gateway on %s:%d\n", c.Host, c.Port)
	server.Start()
}
