package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"qa-service/internal/handlers"
	"qa-service/internal/models"
	"qa-service/internal/repository"
	"qa-service/internal/routes"
	"qa-service/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	db         *gorm.DB
	router     http.Handler
	testServer *httptest.Server
}

func (suite *IntegrationTestSuite) SetupSuite() {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://postgres:password@postgres-test:5432/qa_service_test?sslmode=disable"
	}

	var err error
	suite.db, err = gorm.Open(postgres.Open(testDBURL), &gorm.Config{})
	if err != nil {
		suite.T().Skipf("Skipping tests: cannot connect to test database: %v", err)
		return
	}

	err = suite.db.AutoMigrate(&models.Question{}, &models.Answer{})
	if err != nil {
		suite.T().Fatalf("Failed to migrate test database: %v", err)
	}

	suite.db.Exec("DELETE FROM answers")
	suite.db.Exec("DELETE FROM questions")

	questionRepo := repository.NewQuestionRepository(suite.db)
	answerRepo := repository.NewAnswerRepository(suite.db)
	questionService := services.NewQuestionService(questionRepo)
	answerService := services.NewAnswerService(answerRepo, questionRepo)

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	questionHandler := handlers.NewQuestionHandler(questionService, logger)
	answerHandler := handlers.NewAnswerHandler(answerService, logger)

	router := routes.SetupRoutes(questionHandler, answerHandler, logger)
	suite.router = router
	suite.testServer = httptest.NewServer(router)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.testServer != nil {
		suite.testServer.Close()
	}
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *IntegrationTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM answers")
	suite.db.Exec("DELETE FROM questions")
}

func TestDatabaseConnection(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://postgres:password@postgres-test:5432/qa_service_test?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test migration
	err = db.AutoMigrate(&models.Question{}, &models.Answer{})
	assert.NoError(t, err)

	// Clean up
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Test basic operations
	question := &models.Question{
		Text: "Test question",
	}
	err = db.Create(question).Error
	assert.NoError(t, err)
	assert.NotZero(t, question.ID)

	// Test retrieval
	var retrieved models.Question
	err = db.First(&retrieved, question.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test question", retrieved.Text)

	// Clean up
	db.Delete(question)
}

func (suite *IntegrationTestSuite) TestCreateAndGetQuestion() {
	questionReq := map[string]string{"text": "Test question?"}
	reqBody, _ := json.Marshal(questionReq)

	resp, err := http.Post(suite.testServer.URL+"/api/v1/questions/", "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdQuestion models.Question
	err = json.NewDecoder(resp.Body).Decode(&createdQuestion)
	assert.NoError(suite.T(), err)
	resp.Body.Close()

	assert.Equal(suite.T(), "Test question?", createdQuestion.Text)
	assert.NotZero(suite.T(), createdQuestion.ID)

	resp, err = http.Get(suite.testServer.URL + "/api/v1/questions/" + fmt.Sprintf("%d", createdQuestion.ID))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedQuestion models.Question
	err = json.NewDecoder(resp.Body).Decode(&retrievedQuestion)
	assert.NoError(suite.T(), err)
	resp.Body.Close()

	assert.Equal(suite.T(), createdQuestion.ID, retrievedQuestion.ID)
	assert.Equal(suite.T(), "Test question?", retrievedQuestion.Text)
}

func (suite *IntegrationTestSuite) TestGetAllQuestions() {
	questions := []string{"Question 1?", "Question 2?", "Question 3?"}

	for _, q := range questions {
		req := map[string]string{"text": q}
		reqBody, _ := json.Marshal(req)
		resp, _ := http.Post(suite.testServer.URL+"/api/v1/questions/", "application/json", bytes.NewBuffer(reqBody))
		resp.Body.Close()
	}

	resp, err := http.Get(suite.testServer.URL + "/api/v1/questions/")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedQuestions []models.Question
	err = json.NewDecoder(resp.Body).Decode(&retrievedQuestions)
	assert.NoError(suite.T(), err)
	resp.Body.Close()

	assert.Len(suite.T(), retrievedQuestions, 3)
}

func (suite *IntegrationTestSuite) TestCreateAnswer() {
	questionReq := map[string]string{"text": "Test question for answer?"}
	reqBody, _ := json.Marshal(questionReq)

	resp, err := http.Post(suite.testServer.URL+"/api/v1/questions/", "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(suite.T(), err)

	var question models.Question
	err = json.NewDecoder(resp.Body).Decode(&question)
	assert.NoError(suite.T(), err)
	resp.Body.Close()

	answerReq := map[string]string{
		"user_id": "user123",
		"text":    "This is an answer",
	}
	reqBody, _ = json.Marshal(answerReq)

	resp, err = http.Post(suite.testServer.URL+"/api/v1/questions/"+fmt.Sprintf("%d", question.ID)+"/answers/", "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdAnswer models.Answer
	err = json.NewDecoder(resp.Body).Decode(&createdAnswer)
	assert.NoError(suite.T(), err)
	resp.Body.Close()

	assert.Equal(suite.T(), question.ID, createdAnswer.QuestionID)
	assert.Equal(suite.T(), "user123", createdAnswer.UserID)
	assert.Equal(suite.T(), "This is an answer", createdAnswer.Text)
}

func (suite *IntegrationTestSuite) TestDeleteQuestion() {
	questionReq := map[string]string{"text": "Question to delete"}
	reqBody, _ := json.Marshal(questionReq)

	resp, err := http.Post(suite.testServer.URL+"/api/v1/questions/", "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(suite.T(), err)

	var question models.Question
	err = json.NewDecoder(resp.Body).Decode(&question)
	assert.NoError(suite.T(), err)
	resp.Body.Close()

	req, _ := http.NewRequest("DELETE", suite.testServer.URL+"/api/v1/questions/"+fmt.Sprintf("%d", question.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	resp, err = http.Get(suite.testServer.URL + "/api/v1/questions/" + fmt.Sprintf("%d", question.ID))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
