package models

type TransactionDetail struct {
	ID            string `json:"id" gorm:"type:char(26);primaryKey;autoIncrement:false"`
	TransactionID string `json:"transaction_id" gorm:"type:char(26)"`
	ProductID     string `json:"product_id" gorm:"type:char(26)"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}
