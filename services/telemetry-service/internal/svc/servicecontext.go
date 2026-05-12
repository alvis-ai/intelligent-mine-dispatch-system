package svc

import (
	"context"
	"time"

	"github.com/aicong/mine-dispatch/services/telemetry-service/internal/config"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:        c.RedisAddr,
		Password:    c.RedisPass,
		DB:          c.RedisDB,
		DialTimeout: 5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	return &ServiceContext{Config: c, Redis: rdb}
}
