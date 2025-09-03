package models

type User struct {
	Id       uint    `json:"id" gorm:"primaryKey"`
	Username string  `json:"username" gorm:"unique"`
	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Exp      int     `json:"exp"`
	Classes  []Class `gorm:"many2many:user_classes;" json:"classes"`
	Token    *Token  `json:"token"`
}
