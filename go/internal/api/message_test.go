package api

import (
	"bytes"
	"chat_system/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMessage(t *testing.T) {
	// Set up the Gin engine
	r := gin.Default()

	// Create the mock Redis client
	mockRedis := new(MockRedisClientWrapper)

	// Define the API route
	r.POST("/messages/:application_token/:chat_number", func(c *gin.Context) {
		CreateMessage(c, mockRedis)
	})

	// Test case: Valid message creation
	t.Run("valid message creation", func(t *testing.T) {
		message := models.Message{
			ApplicationToken: "testAppToken",
			ChatNumber:       1,
			Content:          "Test message",
		}

		// Mock Redis behavior
		chatHashKey := "chat#testAppToken-1"
		mockRedis.On("KeyExists", mock.Anything, "chats_messages_count", chatHashKey).Return(true, nil)
		mockRedis.On("IncrementField", mock.Anything, "chats_messages_count", chatHashKey).Return(int64(1), nil)
		mockRedis.On("PushToQueue", mock.Anything, "messages_queue", mock.Anything).Return(nil)

		// Create the HTTP request with valid data
		jsonData, _ := json.Marshal(message)
		req, _ := http.NewRequest(http.MethodPost, "/messages/testAppToken/1", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Message
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, message.ApplicationToken, response.ApplicationToken)
		assert.Equal(t, 1, response.Number) // Asserting incremented message number
		mockRedis.AssertExpectations(t)

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Invalid chat number
	t.Run("invalid chat number", func(t *testing.T) {
		// Create the HTTP request with an invalid chat number
		req, _ := http.NewRequest(http.MethodPost, "/messages/testAppToken/invalid", bytes.NewReader([]byte(`{}`)))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid chat number")

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Chat not found
	t.Run("chat not found", func(t *testing.T) {
		// Mock Redis behavior for non-existing chat
		chatHashKey := "chat#testAppToken-1"
		mockRedis.On("KeyExists", mock.Anything, "chats_messages_count", chatHashKey).Return(false, nil)

		// Create the HTTP request with valid data
		message := models.Message{
			ApplicationToken: "testAppToken",
			ChatNumber:       1,
			Content:          "Test message",
		}
		jsonData, _ := json.Marshal(message)
		req, _ := http.NewRequest(http.MethodPost, "/messages/testAppToken/1", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Chat not found", response["error"])

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Redis error on incrementing the message number
	t.Run("failed to increment message number", func(t *testing.T) {
		// Mock Redis behavior
		chatHashKey := "chat#testAppToken-1"
		mockRedis.On("KeyExists", mock.Anything, "chats_messages_count", chatHashKey).Return(true, nil)
		mockRedis.On("IncrementField", mock.Anything, "chats_messages_count", chatHashKey).Return(int64(0), assert.AnError)

		// Create the HTTP request with valid data
		message := models.Message{
			ApplicationToken: "testAppToken",
			ChatNumber:       1,
			Content:          "Test message",
		}
		jsonData, _ := json.Marshal(message)
		req, _ := http.NewRequest(http.MethodPost, "/messages/testAppToken/1", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to increment message number", response["error"])

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Redis error on pushing to the queue
	t.Run("failed to push message to Redis queue", func(t *testing.T) {
		// Mock Redis behavior
		chatHashKey := "chat#testAppToken-1"
		mockRedis.On("KeyExists", mock.Anything, "chats_messages_count", chatHashKey).Return(true, nil)
		mockRedis.On("IncrementField", mock.Anything, "chats_messages_count", chatHashKey).Return(int64(1), nil)
		mockRedis.On("PushToQueue", mock.Anything, "messages_queue", mock.Anything).Return(assert.AnError)

		// Create the HTTP request with valid data
		message := models.Message{
			ApplicationToken: "testAppToken",
			ChatNumber:       1,
			Content:          "Test message",
		}
		jsonData, _ := json.Marshal(message)
		req, _ := http.NewRequest(http.MethodPost, "/messages/testAppToken/1", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to push message to Redis queue", response["error"])

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})
}
