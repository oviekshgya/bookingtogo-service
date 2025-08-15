package domain

const NATIONALITY = "nationality"

type Nationality struct {
	ID   int    `gorm:"primaryKey;autoIncrement;column:nationality_id"`
	Name string `gorm:"size:50;not null;column:nationality_name"`
	Code string `gorm:"size:2;not null;column:nationality_code"`
}

func (Nationality) TableName() string {
	return NATIONALITY
}
