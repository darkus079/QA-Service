package repository

import (
	"qa-service/internal/models"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(question *models.Question) error {
	return r.db.Create(question).Error
}

func (r *QuestionRepository) GetAll() ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Order("created_at DESC").Find(&questions).Error
	return questions, err
}

func (r *QuestionRepository) GetByID(id uint) (*models.Question, error) {
	var question models.Question
	err := r.db.Preload("Answers").First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *QuestionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Question{}, id).Error
}

func (r *QuestionRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Question{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
