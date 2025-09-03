package models

type Class struct {
	Id              uint           `json:"id" gorm:"primaryKey"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Users           []User         `gorm:"many2many:user_classes;" json:"users"`
	ClassCategoryID uint           `gorm:"not null" json:"class_category_id"`
	ClassCategory   *ClassCategory `gorm:"foreignKey:ClassCategoryID;constraint:OnDelete:SET NULL;references:ID" json:"class_category"`
}
