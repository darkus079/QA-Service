package models

import (
	"time"
)

type Question struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" gorm:"not null" validate:"required,min=1,max=1000"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Answers   []Answer  `json:"answers,omitempty" gorm:"foreignKey:QuestionID"`
}

type CreateQuestionRequest struct {
	Text string `json:"text" validate:"required,min=1,max=1000"`
}
