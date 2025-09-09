package models

type Module struct {
	Id      uint        `json:"id" gorm:"primaryKey"`
	Title   string      `json:"title"`
	Order   int         `json:"order"`
	ClassID uint        `json:"class_id"`
	SubMods []SubModule `json:"submodules"`
}
