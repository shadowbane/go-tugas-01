package models

type Product struct {
	ID          string `json:"id" gorm:"type:char(26);primaryKey;autoIncrement:false"`
	ProductId   string `json:"product_id" gorm:"type:varchar(255);not null"`
	Name        string `json:"name" gorm:"type:varchar(255);null"`
	Description string `json:"description" gorm:"type:varchar(255);null"`
	Price       int    `json:"price" gorm:"type:int;null"`
	Stock       int    `json:"stock" gorm:"type:int;null"`
}
