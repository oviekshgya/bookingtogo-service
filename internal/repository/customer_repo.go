package repository

import (
	"any/bookingtogo-service/internal/domain"
	"any/bookingtogo-service/src/pkg"
	"time"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(db *gorm.DB, customer *domain.Customer) error
	Update(db *gorm.DB, customer *domain.Customer) error
	DeleteCustomer(tx *gorm.DB, customerID uint) error
	GetCustomerByID(tx *gorm.DB, customerID uint) (*domain.Customer, error)
	ListCustomersByNationalityID(tx *gorm.DB, nationalityID uint) ([]domain.Customer, error)
	ListAllCustomers(tx *gorm.DB, page int, pageSize int) (*PaginatedCustomers, error)
}
type CustomerRepositoryImpl struct {
}

func NewUserRepository() CustomerRepository {
	return &CustomerRepositoryImpl{}
}

func (r *CustomerRepositoryImpl) Create(db *gorm.DB, customer *domain.Customer) error {
	return db.Create(customer).Error
}
func (r *CustomerRepositoryImpl) Update(db *gorm.DB, customer *domain.Customer) error {
	var existing domain.Customer
	if err := db.First(&existing, customer.ID).Error; err != nil {
		return err
	}

	data := pkg.UpdateFieldsDynamic(customer)
	// misal ingin simpan updated_at
	data["updated_at"] = time.Now()

	return db.Model(&domain.Customer{}).
		Where("cat_id = ?", customer.ID).
		Updates(data).Error
}

// DeleteCustomer hapus customer berdasarkan ID
func (r *CustomerRepositoryImpl) DeleteCustomer(tx *gorm.DB, customerID uint) error {
	return tx.Delete(&domain.Customer{}, customerID).Error
}

// GetCustomerByID ambil data customer beserta family dan nationality
func (r *CustomerRepositoryImpl) GetCustomerByID(tx *gorm.DB, customerID uint) (*domain.Customer, error) {
	var customer domain.Customer
	err := tx.Preload("FamilyList").Preload("Nationality").
		First(&customer, customerID).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepositoryImpl) ListCustomersByNationalityID(tx *gorm.DB, nationalityID uint) ([]domain.Customer, error) {
	var customers []domain.Customer
	err := tx.Where("nationality_id = ?", nationalityID).
		Preload("FamilyList").Preload("Nationality").
		Order("cat_name asc").
		Find(&customers).Error
	return customers, err
}

type PaginatedCustomers struct {
	Customers []domain.Customer `json:"customers"`
	TotalData int64             `json:"total_data"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
	Page      int               `json:"page"`
}

// ListAllCustomers ambil semua customer dengan pagination
func (r *CustomerRepositoryImpl) ListAllCustomers(tx *gorm.DB, page int, pageSize int) (*PaginatedCustomers, error) {
	var customers []domain.Customer
	var totalData int64

	// Hitung total data
	if err := tx.Model(&domain.Customer{}).Count(&totalData).Error; err != nil {
		return nil, err
	}

	// Hitung offset
	offset := (page - 1) * pageSize

	err := tx.Preload("FamilyList").Preload("Nationality").
		Order("cat_name asc").
		Limit(pageSize).
		Offset(offset).
		Find(&customers).Error
	if err != nil {
		return nil, err
	}

	// Hitung total page
	totalPage := int((totalData + int64(pageSize) - 1) / int64(pageSize)) // ceil division

	return &PaginatedCustomers{
		Customers: customers,
		TotalData: totalData,
		PageSize:  pageSize,
		TotalPage: totalPage,
		Page:      page,
	}, nil
}
