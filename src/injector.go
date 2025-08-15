package src

import (
	"any/bookingtogo-service/internal/handler"
	"any/bookingtogo-service/internal/repository"
	"any/bookingtogo-service/internal/service"
	"any/bookingtogo-service/src/db"
	rds "any/bookingtogo-service/src/redis"

	"github.com/google/wire"
)

func InitializeGlobalConfig() (handler.GlobalConfig, error) {
	wire.Build(
		rds.NewRedisClient,
		db.GetDB,
		repository.NewUserRepository,
		service.NewCustomerService,
		service.NewNasionalityService,
		handler.NewGlobalHandler,
		repository.NewNasionalityRepository,
		service.NewRequestLogService,
		repository.NewRequestLogRepository,
	)
	return nil, nil
}
