package service

import (
	"any/bookingtogo-service/internal/repository"
	"any/bookingtogo-service/src/pkg"

	"gorm.io/gorm"
)

type RequestLogService interface {
	ListAll(page int, pageSize int) (*repository.PaginatedLogs, error)
}

type RequestLogServiceImpl struct {
	DB      *gorm.DB
	LogRepo repository.RequestLogRepository
}

func NewRequestLogService(db *gorm.DB, repo repository.RequestLogRepository) RequestLogService {
	return &RequestLogServiceImpl{
		DB:      db,
		LogRepo: repo,
	}
}

// ListAll ambil semua log dengan pagination
func (s *RequestLogServiceImpl) ListAll(page int, pageSize int) (*repository.PaginatedLogs, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		return s.LogRepo.ListLogs(tx, page, pageSize)
	})
	if err != nil {
		return nil, err
	}
	return result.(*repository.PaginatedLogs), nil
}
