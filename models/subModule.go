package models

type SubModule struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	ModuleID    uint   `json:"module_id"`
}
