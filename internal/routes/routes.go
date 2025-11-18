package routes

import (
	"log"
	"net/http"
	"qa-service/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(questionHandler *handlers.QuestionHandler, answerHandler *handlers.AnswerHandler, logger *log.Logger) *mux.Router {
	router := mux.NewRouter()

	router.Use(loggingMiddleware(logger))

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/questions/", questionHandler.GetQuestions).Methods("GET")
	api.HandleFunc("/questions/", questionHandler.CreateQuestion).Methods("POST")
	api.HandleFunc("/questions/{id:[0-9]+}", questionHandler.GetQuestion).Methods("GET")
	api.HandleFunc("/questions/{id:[0-9]+}", questionHandler.DeleteQuestion).Methods("DELETE")

	api.HandleFunc("/questions/{id:[0-9]+}/answers/", answerHandler.CreateAnswer).Methods("POST")
	api.HandleFunc("/answers/{id:[0-9]+}", answerHandler.GetAnswer).Methods("GET")
	api.HandleFunc("/answers/{id:[0-9]+}", answerHandler.DeleteAnswer).Methods("DELETE")

	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	return router
}

func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}
