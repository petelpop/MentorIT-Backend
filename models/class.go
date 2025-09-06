package models

type Class struct {
	Id              uint           `json:"id" gorm:"primaryKey"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Thumbnail       string         `json:"thumbnail" form:"thumbnail"`
	Trailer         string         `json:"trailer"`
	Price           int            `json:"price"`
	FormattedPrice  string         `json:"formatted_price"`
	CategoryName    string         `json:"category_name"`
	ClassCategoryID uint           `gorm:"not null" json:"class_category_id"`
	ClassCategory   *ClassCategory `gorm:"constraint:OnDelete:SET NULL;" json:"class_category"`
	Users           []User         `gorm:"many2many:user_classes;" json:"users"`
	// Mentor          []Mentor       `gorm:"many2many:mentor_classes;" json:"mentor"`
}
