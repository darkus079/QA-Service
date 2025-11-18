package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"qa-service/internal/database"
	"qa-service/internal/handlers"
	"qa-service/internal/repository"
	"qa-service/internal/routes"
	"qa-service/internal/services"
	"syscall"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "[QA-SERVICE] ", log.LstdFlags)

	logger.Println("Initializing database...")
	if err := database.InitDB(); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			logger.Printf("Error closing database: %v", err)
		}
	}()

	questionRepo := repository.NewQuestionRepository(database.GetDB())
	answerRepo := repository.NewAnswerRepository(database.GetDB())

	questionService := services.NewQuestionService(questionRepo)
	answerService := services.NewAnswerService(answerRepo, questionRepo)

	questionHandler := handlers.NewQuestionHandler(questionService, logger)
	answerHandler := handlers.NewAnswerHandler(answerService, logger)

	router := routes.SetupRoutes(questionHandler, answerHandler, logger)

	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	shutdownTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdownChan := make(chan struct{})
	go func() {
		if err := server.Shutdown(ctx); err != nil {
			logger.Printf("Server shutdown error: %v", err)
		}
		close(shutdownChan)
	}()

	select {
	case <-shutdownChan:
		logger.Println("Server shutdown complete")
	case <-time.After(shutdownTimeout):
		logger.Println("Server shutdown timeout")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
