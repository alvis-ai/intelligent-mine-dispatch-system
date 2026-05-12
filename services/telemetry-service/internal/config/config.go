package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	RedisAddr string
	RedisPass string
	RedisDB   int
}
