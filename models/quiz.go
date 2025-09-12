package models

type Quiz struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title"`     
	ModuleID  uint           `json:"module_id"` 
	Questions []QuizQuestion `json:"questions" gorm:"foreignKey:QuizID"`
}
