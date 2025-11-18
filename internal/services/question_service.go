package services

import (
	"errors"
	"qa-service/internal/models"
	"qa-service/internal/repository"
)

type QuestionService struct {
	questionRepo *repository.QuestionRepository
}

func NewQuestionService(questionRepo *repository.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
	}
}

func (s *QuestionService) CreateQuestion(req *models.CreateQuestionRequest) (*models.Question, error) {
	if req.Text == "" {
		return nil, errors.New("question text cannot be empty")
	}

	question := &models.Question{
		Text: req.Text,
	}

	err := s.questionRepo.Create(question)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (s *QuestionService) GetAllQuestions() ([]models.Question, error) {
	return s.questionRepo.GetAll()
}

func (s *QuestionService) GetQuestionByID(id uint) (*models.Question, error) {
	return s.questionRepo.GetByID(id)
}

func (s *QuestionService) DeleteQuestion(id uint) error {
	exists, err := s.questionRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("question not found")
	}

	return s.questionRepo.Delete(id)
}
