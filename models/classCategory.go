package models

type ClassCategory struct {
	Id          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"unique"`
	Description string  `json:"description"`
	Icon        string  `json:"icon" form:"icon"`
	Classes     []Class `json:"classes"`
}
