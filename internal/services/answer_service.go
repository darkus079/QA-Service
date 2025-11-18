package services

import (
	"errors"
	"qa-service/internal/models"
	"qa-service/internal/repository"
)

type AnswerService struct {
	answerRepo   *repository.AnswerRepository
	questionRepo *repository.QuestionRepository
}

func NewAnswerService(answerRepo *repository.AnswerRepository, questionRepo *repository.QuestionRepository) *AnswerService {
	return &AnswerService{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
	}
}

func (s *AnswerService) CreateAnswer(questionID uint, req *models.CreateAnswerRequest) (*models.Answer, error) {
	if req.Text == "" {
		return nil, errors.New("answer text cannot be empty")
	}
	if req.UserID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	exists, err := s.questionRepo.Exists(questionID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("question not found")
	}

	answer := &models.Answer{
		QuestionID: questionID,
		UserID:     req.UserID,
		Text:       req.Text,
	}

	err = s.answerRepo.Create(answer)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (s *AnswerService) GetAnswerByID(id uint) (*models.Answer, error) {
	return s.answerRepo.GetByID(id)
}

func (s *AnswerService) DeleteAnswer(id uint) error {
	exists, err := s.answerRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("answer not found")
	}

	return s.answerRepo.Delete(id)
}
