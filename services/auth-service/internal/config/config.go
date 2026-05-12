package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	JwtSecret   string
	JwtExpire   int64
	UserRpc     zrpc.RpcClientConf
	PostgresDSN string
	RedisAddr   string
	RedisPass   string
}
