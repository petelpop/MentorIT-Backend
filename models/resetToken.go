package models

import "time"

type ResetToken struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token" gorm:"unique"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID"`
}
