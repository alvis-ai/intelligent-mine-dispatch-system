package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/aicong/mine-dispatch/services/auth-service/internal/config"
	"github.com/aicong/mine-dispatch/services/auth-service/internal/logic"
	"github.com/aicong/mine-dispatch/services/auth-service/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/auth.yaml", "config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.SetLevel(logx.InfoLevel)

	ctx := svc.NewServiceContext(c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/auth/login",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req logic.LoginRequest
			json.Unmarshal(body, &req)
			l := logic.NewLoginLogic(r.Context(), ctx)
			resp, _ := l.Login(&req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/auth/validate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req logic.ValidateRequest
			json.Unmarshal(body, &req)
			l := logic.NewValidateLogic(r.Context(), ctx)
			resp, _ := l.Validate(&req)
			data, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		},
	})

	fmt.Printf("Starting auth-service on %s:%d\n", c.Host, c.Port)
	server.Start()
}
