package domain

import "time"

const CUSTOMER = "customer"

type Customer struct {
	ID            int         `gorm:"primaryKey;autoIncrement;column:cat_id"`
	Name          string      `gorm:"size:50;not null;column:cat_name"`
	DOB           time.Time   `gorm:"not null;column:cat_dob"`
	PhoneNumber   string      `gorm:"size:20;not null;column:cat_phoneNum"`
	Email         string      `gorm:"size:50;not null;column:cat_email"`
	NationalityID int         `gorm:"not null;column:nationality_id"`
	Nationality   Nationality `gorm:"foreignKey:NationalityID"`
	FamilyList    []Family    `gorm:"foreignKey:CustomerID"`
}

func (Customer) TableName() string {
	return CUSTOMER
}
