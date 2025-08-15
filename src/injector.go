package src

import (
	"any/bookingtogo-service/internal/handler"
	"any/bookingtogo-service/internal/repository"
	"any/bookingtogo-service/internal/service"
	"any/bookingtogo-service/src/db"
	rds "any/bookingtogo-service/src/redis"

	"github.com/google/wire"
)

func InitializeNasionalityControllers() (handler.NasionalityHandler, error) {
	wire.Build(
		handler.NewNasionalityHandler,
		rds.NewRedisClient,
		db.GetDB,
		repository.NewNasionalityRepository,
		service.NewNasionalityService,
	)
	return nil, nil
}

func InitializeCustomerControllers() (handler.CustomerHandler, error) {
	wire.Build(
		handler.NewCustomerHandler,
		rds.NewRedisClient,
		db.GetDB,
		repository.NewUserRepository,
		service.NewCustomerService,
	)
	return nil, nil
}
