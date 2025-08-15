package service

import (
	"any/bookingtogo-service/internal/repository"
	"any/bookingtogo-service/src/pkg"
	"any/bookingtogo-service/src/redis"

	"gorm.io/gorm"
)

type NasionalityServiceImpl struct {
	DB              *gorm.DB
	NasionalityRepo repository.NasionalityRepository
	Redis           *redis.RedisClient
}

type NasionalityService interface {
	GetAll() (interface{}, error)
	GetById(id int) (interface{}, error)
}

func NewNasionalityService(db *gorm.DB, redis *redis.RedisClient, repo repository.NasionalityRepository) NasionalityService {
	return &NasionalityServiceImpl{
		Redis:           redis,
		NasionalityRepo: repo,
		DB:              db,
	}
}

func (s *NasionalityServiceImpl) GetAll() (interface{}, error) {
	return pkg.WithTransaction(s.DB, func(tz *gorm.DB) (interface{}, error) {
		data, err := s.NasionalityRepo.ListNationalities(tz)
		if err != nil {
			return nil, err
		}
		return data, nil
	})
}

func (s *NasionalityServiceImpl) GetById(id int) (interface{}, error) {
	return pkg.WithTransaction(s.DB, func(tz *gorm.DB) (interface{}, error) {
		data, err := s.NasionalityRepo.GetNationalityByID(tz, uint(id))
		if err != nil {
			return nil, err
		}
		return data, nil
	})
}
