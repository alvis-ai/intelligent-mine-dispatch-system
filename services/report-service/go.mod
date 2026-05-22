module github.com/aicong/mine-dispatch/services/report-service

go 1.25.3

replace github.com/aicong/mine-dispatch/proto => ../../proto

replace github.com/aicong/mine-dispatch/pkg => ../../pkg

require (
	github.com/aicong/mine-dispatch/pkg v0.0.0-00010101000000-000000000000
	github.com/aicong/mine-dispatch/proto v0.0.0-00010101000000-000000000000
	github.com/zeromicro/go-zero v1.10.1
	google.golang.org/grpc v1.81.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)
