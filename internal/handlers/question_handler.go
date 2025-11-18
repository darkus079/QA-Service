package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"qa-service/internal/models"
	"qa-service/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

type QuestionHandler struct {
	questionService *services.QuestionService
	logger          *log.Logger
}

func NewQuestionHandler(questionService *services.QuestionService, logger *log.Logger) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
		logger:          logger,
	}
}

func (h *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling GET /questions/")

	questions, err := h.questionService.GetAllQuestions()
	if err != nil {
		h.logger.Printf("Error getting questions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		h.logger.Printf("Error encoding questions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *QuestionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling POST /questions/")

	var req models.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.CreateQuestion(&req)
	if err != nil {
		h.logger.Printf("Error creating question: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(question); err != nil {
		h.logger.Printf("Error encoding question: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *QuestionHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Printf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Handling GET /questions/%d", id)

	question, err := h.questionService.GetQuestionByID(uint(id))
	if err != nil {
		h.logger.Printf("Error getting question: %v", err)
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(question); err != nil {
		h.logger.Printf("Error encoding question: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *QuestionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Printf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Handling DELETE /questions/%d", id)

	err = h.questionService.DeleteQuestion(uint(id))
	if err != nil {
		h.logger.Printf("Error deleting question: %v", err)
		if err.Error() == "question not found" {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
