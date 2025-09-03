package models

import (
	"time"
)

type Token struct {
	Id           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `gorm:"not null"`
	User         *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user_data"`
	AccessToken  string    `json:"access_token" gorm:"unique;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"unique;not null"`
	ExpiresAt    time.Time `json:"expires_at"`
}
