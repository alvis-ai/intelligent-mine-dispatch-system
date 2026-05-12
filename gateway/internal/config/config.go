package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	AuthRpc      struct { Etcd struct { Hosts []string `json:"hosts"`; Key string `json:"key"` } `json:"etcd"`; Endpoints []string `json:"endpoints"` } `json:"auth_rpc"`
	UserRpc      struct { Etcd struct { Hosts []string `json:"hosts"`; Key string `json:"key"` } `json:"etcd"`; Endpoints []string `json:"endpoints"` } `json:"user_rpc"`
	VehicleRpc   struct { Etcd struct { Hosts []string `json:"hosts"`; Key string `json:"key"` } `json:"etcd"`; Endpoints []string `json:"endpoints"` } `json:"vehicle_rpc"`
	TelemetryRpc struct { Etcd struct { Hosts []string `json:"hosts"`; Key string `json:"key"` } `json:"etcd"`; Endpoints []string `json:"endpoints"` } `json:"telemetry_rpc"`
	DispatchRpc  struct { Etcd struct { Hosts []string `json:"hosts"`; Key string `json:"key"` } `json:"etcd"`; Endpoints []string `json:"endpoints"` } `json:"dispatch_rpc"`
}
