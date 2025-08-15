package handler

import (
	"any/bookingtogo-service/internal/repository"
	"any/bookingtogo-service/internal/service"
	"any/bookingtogo-service/src/redis"

	"gorm.io/gorm"
)

type GlobalConfig interface {
	GetConnectionRedis() *redis.RedisClient
	GetConnectionDB() *gorm.DB
	ServiceCustomer() service.CustomerService
	ServiceNas() service.NasionalityService
	RepositoryCustomer() repository.CustomerRepository
	RepositoryNas() repository.NasionalityRepository
	LogService() service.RequestLogService
	LogRepositorys() repository.RequestLogRepository
}

type GlobalHandlerImpl struct {
	DB              *gorm.DB
	Redis           *redis.RedisClient
	Service         service.CustomerService
	NasionalityRepo repository.NasionalityRepository
	CustomerRepo    repository.CustomerRepository
	ServiceNasi     service.NasionalityService
	LogSerrvice     service.RequestLogService
	LogRepository   repository.RequestLogRepository
}

func NewGlobalHandler(db *gorm.DB, rdis *redis.RedisClient, serviceCus service.CustomerService, serviceNas service.NasionalityService, repoCust repository.CustomerRepository, repoNAs repository.NasionalityRepository, lg service.RequestLogService, lgRp repository.RequestLogRepository) GlobalConfig {
	return &GlobalHandlerImpl{
		DB:              db,
		Redis:           rdis,
		Service:         serviceCus,
		NasionalityRepo: repoNAs,
		CustomerRepo:    repoCust,
		ServiceNasi:     serviceNas,
		LogSerrvice:     lg,
		LogRepository:   lgRp,
	}
}

func (g *GlobalHandlerImpl) GetConnectionRedis() *redis.RedisClient {
	return g.Redis
}

func (g *GlobalHandlerImpl) GetConnectionDB() *gorm.DB {
	return g.DB
}

func (g *GlobalHandlerImpl) ServiceCustomer() service.CustomerService {
	return g.Service
}

func (g *GlobalHandlerImpl) ServiceNas() service.NasionalityService {
	return g.ServiceNasi
}

func (g *GlobalHandlerImpl) RepositoryCustomer() repository.CustomerRepository {
	return g.CustomerRepo
}

func (g *GlobalHandlerImpl) RepositoryNas() repository.NasionalityRepository {
	return g.NasionalityRepo
}

func (g *GlobalHandlerImpl) LogService() service.RequestLogService {
	return g.LogSerrvice
}

func (g *GlobalHandlerImpl) LogRepositorys() repository.RequestLogRepository {
	return g.LogRepository
}
