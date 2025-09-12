package models

type QuizQuestion struct {
    Id       uint   `json:"id" gorm:"primaryKey"`
    QuizID   uint   `json:"quiz_id"`   
    Question string `json:"question"`   
    Options  string `json:"options"`  
    Answer   string `json:"answer"`     
    Order    int    `json:"order"` 
}
