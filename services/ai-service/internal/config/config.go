package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	PostgresDSN  string
	RedisAddr    string
	RedisPass    string
	RouteSvcAddr string
}
