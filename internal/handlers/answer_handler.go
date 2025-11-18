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

type AnswerHandler struct {
	answerService *services.AnswerService
	logger        *log.Logger
}

func NewAnswerHandler(answerService *services.AnswerService, logger *log.Logger) *AnswerHandler {
	return &AnswerHandler{
		answerService: answerService,
		logger:        logger,
	}
}

func (h *AnswerHandler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Question ID not provided", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		h.logger.Printf("Invalid question ID: %v", err)
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Handling POST /questions/%d/answers/", questionID)

	var req models.CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	answer, err := h.answerService.CreateAnswer(uint(questionID), &req)
	if err != nil {
		h.logger.Printf("Error creating answer: %v", err)
		if err.Error() == "question not found" {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		h.logger.Printf("Error encoding answer: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *AnswerHandler) GetAnswer(w http.ResponseWriter, r *http.Request) {
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

	h.logger.Printf("Handling GET /answers/%d", id)

	answer, err := h.answerService.GetAnswerByID(uint(id))
	if err != nil {
		h.logger.Printf("Error getting answer: %v", err)
		http.Error(w, "Answer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		h.logger.Printf("Error encoding answer: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *AnswerHandler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
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

	h.logger.Printf("Handling DELETE /answers/%d", id)

	err = h.answerService.DeleteAnswer(uint(id))
	if err != nil {
		h.logger.Printf("Error deleting answer: %v", err)
		if err.Error() == "answer not found" {
			http.Error(w, "Answer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
