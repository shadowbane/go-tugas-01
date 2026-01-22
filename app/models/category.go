package models

// Shadowbane's here. gorm for future update,
// in case the project (idk, depends on the mentor, tho) needs this
// just leave this here, it doesn't have any effect tho

type Category struct {
	ID          string `json:"id" gorm:"type:char(26);primaryKey;autoIncrement:false"`
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:varchar(255);null"`
}
