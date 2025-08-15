package domain

const CUSTOMER = "customer"

type Customer struct {
	ID            int         `gorm:"primaryKey;autoIncrement;column:cat_id"`
	Name          string      `gorm:"size:50;not null;column:cat_name" json:"name"`
	DOB           string      `gorm:"not null;column:cat_dob" json:"DOB"`
	PhoneNumber   string      `gorm:"size:20;not null;column:cat_phoneNum" json:"phoneNumber"`
	Email         string      `gorm:"size:50;not null;column:cat_email" json:"email"`
	NationalityID int         `gorm:"not null;column:nationality_id" json:"nationalityId"`
	Nationality   Nationality `gorm:"foreignKey:NationalityID" json:"nationality"`
	FamilyList    []Family    `gorm:"foreignKey:CustomerID" json:"familyList"`
}

func (Customer) TableName() string {
	return CUSTOMER
}
