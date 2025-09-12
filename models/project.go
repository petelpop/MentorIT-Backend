package models

type FinalProject struct {
    Id          uint   `json:"id" gorm:"primaryKey"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Submission  string `json:"submission" form:"submission"`
    ClassID     uint   `json:"class_id"`
    UserID      uint   `json:"user_id"`
    ReviewedBy  uint   `json:"reviewed_by"`
    Grade       string `json:"grade"`
    Feedback    string `json:"feedback"`
}
