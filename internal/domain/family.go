package domain

const FAMILY = "family_list"

type Family struct {
	ID         int    `gorm:"primaryKey;autoIncrement;column:fl_id"`
	CustomerID int    `gorm:"not null;column:cat_id"`
	Relation   string `gorm:"size:50;not null;column:fl_relation"`
	Name       string `gorm:"size:50;not null;column:fl_name" json:"name"`
	DOB        string `gorm:"size:50;not null;column:fl_dob" json:"DOB"` // bisa string atau time.Time
}

func (Family) TableName() string {
	return FAMILY
}
