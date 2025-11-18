package repository

import (
	"qa-service/internal/models"

	"gorm.io/gorm"
)

type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Create(answer *models.Answer) error {
	return r.db.Create(answer).Error
}

func (r *AnswerRepository) GetByID(id uint) (*models.Answer, error) {
	var answer models.Answer
	err := r.db.Preload("Question").First(&answer, id).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *AnswerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Answer{}, id).Error
}

func (r *AnswerRepository) GetByQuestionID(questionID uint) ([]models.Answer, error) {
	var answers []models.Answer
	err := r.db.Where("question_id = ?", questionID).Order("created_at ASC").Find(&answers).Error
	return answers, err
}

func (r *AnswerRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Answer{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
