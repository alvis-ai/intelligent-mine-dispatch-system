package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	PostgresDSN string
	RedisAddr   string
	RedisPass   string
}

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitPostgres(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}

func InitRedis(addr, pass string) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
}

func Ping(ctx context.Context) error {
	if DB != nil {
		sqlDB, _ := DB.DB()
		if err := sqlDB.PingContext(ctx); err != nil {
			return err
		}
	}
	if Redis != nil {
		return Redis.Ping(ctx).Err()
	}
	return nil
}
