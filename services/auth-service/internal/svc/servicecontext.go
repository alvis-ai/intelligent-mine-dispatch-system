package svc

import (
	"github.com/aicong/mine-dispatch/services/auth-service/internal/config"
	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserModel struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement:false"`
	Username string `gorm:"size:64;uniqueIndex"`
	Password string `gorm:"size:255"`
}

func (UserModel) TableName() string { return "users" }

type ServiceContext struct {
	Config   config.Config
	UserRpc  userv1.UserServiceClient
	DB       *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.UserRpc)
	db, err := gorm.Open(postgres.Open(c.PostgresDSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return &ServiceContext{
		Config:  c,
		UserRpc: userv1.NewUserServiceClient(conn.Conn()),
		DB:      db,
	}
}
