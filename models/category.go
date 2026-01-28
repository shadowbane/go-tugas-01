package models

type Category struct {
	ID          string `json:"id" gorm:"type:char(26);primaryKey;autoIncrement:false"`
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:varchar(255);null"`
}
