package models

type ProjectPage struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	ModuleID    uint   `json:"module_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Guide       string `json:"guide"`
}
