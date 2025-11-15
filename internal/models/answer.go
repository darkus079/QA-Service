package models

import (
	"time"
)

type Answer struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	UserID     string    `json:"user_id" gorm:"not null" validate:"required"`
	Text       string    `json:"text" gorm:"not null" validate:"required,min=1,max=2000"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	Question   Question  `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}

type CreateAnswerRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Text   string `json:"text" validate:"required,min=1,max=2000"`
}
