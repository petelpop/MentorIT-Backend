package models

import "time"

type Transaction struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	OrderID   string `json:"order_id" gorm:"unique"`
	UserID    uint   `json:"user_id"`
	ClassID   uint   `json:"class_id"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}