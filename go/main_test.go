package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

var mockRedisClient redismock.ClientMock

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/applications/:application_token/chats", createChat)
	router.POST("/applications/:application_token/chats/:chat_number/messages", createMessage)
	return router
}

func TestCreateChatSuccess(t *testing.T) {
	r := setupRouter()

	// Set up mock redis client
	db, mock := redismock.NewClientMock()
	redis_client = db
	mockRedisClient = mock

	// Mock Redis behavior
	mock.ExpectHExists("applications_chats_count", "app#test_token").SetVal(true)
	mock.ExpectHIncrBy("applications_chats_count", "app#test_token", 1).SetVal(1)
	mock.ExpectRPush("chats_queue", []byte(`{"number":1,"application_token":"test_token","title":"Test Chat"}`)).SetVal(1)

	// Create a new chat
	newChat := map[string]interface{}{
		"title": "Test Chat",
	}
	jsonValue, _ := json.Marshal(newChat)
	req, _ := http.NewRequest("POST", "/applications/test_token/chats", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Simulate HTTP request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusCreated, w.Code)
	var createdChat chat
	json.Unmarshal(w.Body.Bytes(), &createdChat)
	assert.Equal(t, 1, createdChat.Number)
	assert.Equal(t, "test_token", createdChat.ApplicationToken)
	assert.Equal(t, "Test Chat", createdChat.Title)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateChatAppNotFound(t *testing.T) {
	r := setupRouter()

	// Set up mock redis client
	db, mock := redismock.NewClientMock()
	redis_client = db
	mockRedisClient = mock

	// Mock Redis behavior (application doesn't exist)
	mock.ExpectHExists("applications_chats_count", "app#missing_token").SetVal(false)

	// Create a new chat with missing application token
	newChat := map[string]interface{}{
		"title": "Test Chat",
	}
	jsonValue, _ := json.Marshal(newChat)
	req, _ := http.NewRequest("POST", "/applications/missing_token/chats", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Simulate HTTP request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateMessageSuccess(t *testing.T) {
	r := setupRouter()

	// Set up mock redis client
	db, mock := redismock.NewClientMock()
	redis_client = db
	mockRedisClient = mock

	// Mock Redis behavior
	mock.ExpectHExists("chats_messages_count", "chat#test_token-1").SetVal(true)
	mock.ExpectHIncrBy("chats_messages_count", "chat#test_token-1", 1).SetVal(1)
	mock.ExpectRPush("messages_queue", []byte(`{"number":1,"chat_number":1,"application_token":"test_token","content":"Test Message"}`)).SetVal(1)

	// Create a new message
	newMessage := map[string]interface{}{
		"content": "Test Message",
	}
	jsonValue, _ := json.Marshal(newMessage)
	req, _ := http.NewRequest("POST", "/applications/test_token/chats/1/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Simulate HTTP request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusCreated, w.Code)
	var createdMessage message
	json.Unmarshal(w.Body.Bytes(), &createdMessage)
	assert.Equal(t, 1, createdMessage.Number)
	assert.Equal(t, "test_token", createdMessage.ApplicationToken)
	assert.Equal(t, "Test Message", createdMessage.Content)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateMessageChatNotFound(t *testing.T) {
	r := setupRouter()

	// Set up mock redis client
	db, mock := redismock.NewClientMock()
	redis_client = db
	mockRedisClient = mock

	// Mock Redis behavior (chat doesn't exist)
	mock.ExpectHExists("chats_messages_count", "chat#missing_token-1").SetVal(false)

	// Create a new message with missing chat
	newMessage := map[string]interface{}{
		"content": "Test Message",
	}
	jsonValue, _ := json.Marshal(newMessage)
	req, _ := http.NewRequest("POST", "/applications/missing_token/chats/1/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Simulate HTTP request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
