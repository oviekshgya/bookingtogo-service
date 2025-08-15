package repository

import (
	"any/bookingtogo-service/internal/domain"
	"any/bookingtogo-service/src/pkg"
	"time"

	"gorm.io/gorm"
)

type NasionalityRepository interface {
	CreateNationality(tx *gorm.DB, nationality *domain.Nationality) error
	UpdateNationality(tx *gorm.DB, nationality *domain.Nationality) error
	DeleteNationality(tx *gorm.DB, nationalityID uint) error
	GetNationalityByID(tx *gorm.DB, nationalityID uint) (*domain.Nationality, error)
	ListNationalities(tx *gorm.DB) ([]domain.Nationality, error)
}

type NationalityRepositoryImpl struct{}

func NewNasionalityRepository() NasionalityRepository {
	return &NationalityRepositoryImpl{}
}

func (r *NationalityRepositoryImpl) CreateNationality(tx *gorm.DB, nationality *domain.Nationality) error {
	return tx.Create(nationality).Error
}

func (r *NationalityRepositoryImpl) UpdateNationality(tx *gorm.DB, nationality *domain.Nationality) error {
	var existing domain.Nationality
	if err := tx.First(&existing, nationality.ID).Error; err != nil {
		return err
	}

	data := pkg.UpdateFieldsDynamic(nationality)
	data["updated_at"] = time.Now() // optional kalau pakai timestamps

	return tx.Model(&domain.Nationality{}).
		Where("nationality_id = ?", nationality.ID).
		Updates(data).Error
}

func (r *NationalityRepositoryImpl) DeleteNationality(tx *gorm.DB, nationalityID uint) error {
	return tx.Delete(&domain.Nationality{}, nationalityID).Error
}

func (r *NationalityRepositoryImpl) GetNationalityByID(tx *gorm.DB, nationalityID uint) (*domain.Nationality, error) {
	var nationality domain.Nationality
	err := tx.First(&nationality, nationalityID).Error
	if err != nil {
		return nil, err
	}
	return &nationality, nil
}

func (r *NationalityRepositoryImpl) ListNationalities(tx *gorm.DB) ([]domain.Nationality, error) {
	var nationalities []domain.Nationality
	err := tx.Order("nationality_name asc").Find(&nationalities).Error
	return nationalities, err
}
