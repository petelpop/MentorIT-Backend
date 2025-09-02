package models

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Users       []User `gorm:"many2many:user_classes;" json:"users"`
}
