package models

type ModuleItem struct {
	Id       uint   `json:"id" gorm:"primaryKey"`
	ModuleID uint   `json:"module_id"`
	ItemType string `json:"item_type"`
	ItemID   uint   `json:"item_id"`
	Order    int    `json:"order"`
}
