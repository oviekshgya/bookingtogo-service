package service

import (
	"any/bookingtogo-service/internal/domain"
	"any/bookingtogo-service/internal/repository"
	"any/bookingtogo-service/src/pkg"
	"any/bookingtogo-service/src/redis"

	"gorm.io/gorm"
)

type CustomerServiceImpl struct {
	DB           *gorm.DB
	CustomerRepo repository.CustomerRepository
	Redis        *redis.RedisClient
}

type CustomerService interface {
	Create(customer *domain.Customer) (*domain.Customer, error)
	Update(customer *domain.Customer) (*domain.Customer, error)
	Delete(id uint) error
	GetById(id uint) (*domain.Customer, error)
	ListByNationalityID(nationalityID uint) ([]domain.Customer, error)
}

func NewCustomerService(db *gorm.DB, redis *redis.RedisClient, repo repository.CustomerRepository) CustomerService {
	return &CustomerServiceImpl{
		DB:           db,
		Redis:        redis,
		CustomerRepo: repo,
	}
}

// Create customer baru
func (s *CustomerServiceImpl) Create(customer *domain.Customer) (*domain.Customer, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		if err := s.CustomerRepo.Create(tx, customer); err != nil {
			return nil, err
		}
		return customer, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*domain.Customer), nil
}

// Update customer
func (s *CustomerServiceImpl) Update(customer *domain.Customer) (*domain.Customer, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		if err := s.CustomerRepo.Update(tx, customer); err != nil {
			return nil, err
		}
		return customer, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*domain.Customer), nil
}

// Delete customer berdasarkan ID
func (s *CustomerServiceImpl) Delete(id uint) error {
	_, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		return nil, s.CustomerRepo.DeleteCustomer(tx, id)
	})
	return err
}

// Get customer berdasarkan ID
func (s *CustomerServiceImpl) GetById(id uint) (*domain.Customer, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		return s.CustomerRepo.GetCustomerByID(tx, id)
	})
	if err != nil {
		return nil, err
	}
	return result.(*domain.Customer), nil
}

// List customer berdasarkan nationality
func (s *CustomerServiceImpl) ListByNationalityID(nationalityID uint) ([]domain.Customer, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		return s.CustomerRepo.ListCustomersByNationalityID(tx, nationalityID)
	})
	if err != nil {
		return nil, err
	}
	return result.([]domain.Customer), nil
}
