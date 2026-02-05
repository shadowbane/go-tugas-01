package models

import "time"

type Transaction struct {
	ID          string              `json:"id" gorm:"type:char(26);primaryKey;autoIncrement:false"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}
