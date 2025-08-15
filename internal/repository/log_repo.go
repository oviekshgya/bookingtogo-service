package repository

import (
	"any/bookingtogo-service/internal/domain"
	"any/bookingtogo-service/src/pkg"
	"time"

	"gorm.io/gorm"
)

type RequestLogRepository interface {
	Create(db *gorm.DB, log *domain.RequestLog) error
	Update(db *gorm.DB, log *domain.RequestLog) error
	DeleteLog(tx *gorm.DB, logID uint) error
	GetLogByID(tx *gorm.DB, logID uint) (*domain.RequestLog, error)
	ListLogs(tx *gorm.DB, page int, pageSize int) (*PaginatedLogs, error)
}

type RequestLogRepositoryImpl struct{}

func NewRequestLogRepository() RequestLogRepository {
	return &RequestLogRepositoryImpl{}
}

func (r *RequestLogRepositoryImpl) Create(db *gorm.DB, log *domain.RequestLog) error {
	return db.Create(log).Error
}

func (r *RequestLogRepositoryImpl) Update(db *gorm.DB, log *domain.RequestLog) error {
	var existing domain.RequestLog
	if err := db.First(&existing, log.ID).Error; err != nil {
		return err
	}

	data := pkg.UpdateFieldsDynamic(log)
	data["updated_at"] = time.Now()

	return db.Model(&domain.RequestLog{}).
		Where("id = ?", log.ID).
		Updates(data).Error
}

func (r *RequestLogRepositoryImpl) DeleteLog(tx *gorm.DB, logID uint) error {
	return tx.Delete(&domain.RequestLog{}, logID).Error
}

func (r *RequestLogRepositoryImpl) GetLogByID(tx *gorm.DB, logID uint) (*domain.RequestLog, error) {
	var log domain.RequestLog
	err := tx.First(&log, logID).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

type PaginatedLogs struct {
	Logs      []domain.RequestLog `json:"logs"`
	TotalData int64               `json:"total_data"`
	PageSize  int                 `json:"page_size"`
	TotalPage int                 `json:"total_page"`
	Page      int                 `json:"page"`
}

func (r *RequestLogRepositoryImpl) ListLogs(tx *gorm.DB, page int, pageSize int) (*PaginatedLogs, error) {
	var logs []domain.RequestLog
	var totalData int64

	// Hitung total data
	if err := tx.Model(&domain.RequestLog{}).Count(&totalData).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize

	err := tx.Order("timestamp desc").
		Limit(pageSize).
		Offset(offset).
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	totalPage := int((totalData + int64(pageSize) - 1) / int64(pageSize))

	return &PaginatedLogs{
		Logs:      logs,
		TotalData: totalData,
		PageSize:  pageSize,
		TotalPage: totalPage,
		Page:      page,
	}, nil
}
