package api

import (
	"bytes"
	"chat_system/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateChat(t *testing.T) {
	// Set up the Gin engine
	r := gin.Default()

	// Create the mock Redis client
	mockRedis := new(MockRedisClientWrapper)

	// Define the API route
	r.POST("/chats/:application_token", func(c *gin.Context) {
		CreateChat(c, mockRedis)
	})

	// Test case: Valid chat creation
	t.Run("valid chat creation", func(t *testing.T) {
		chat := models.Chat{
			ApplicationToken: "testAppToken",
			Number:           1,
			Title:            "testTitle",
		}

		// Mock Redis behavior
		mockRedis.On("KeyExists", mock.Anything, "applications_chats_count", "app#testAppToken").Return(true, nil)
		mockRedis.On("IncrementField", mock.Anything, "applications_chats_count", "app#testAppToken").Return(int64(1), nil)
		mockRedis.On("PushToQueue", mock.Anything, "chats_queue", chat).Return(nil)

		// Create the HTTP request with valid data
		jsonData, _ := json.Marshal(chat)
		req, _ := http.NewRequest(http.MethodPost, "/chats/testAppToken", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Chat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, chat.ApplicationToken, response.ApplicationToken)
		mockRedis.AssertExpectations(t)

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Invalid JSON data
	t.Run("invalid JSON", func(t *testing.T) {
		// Create the HTTP request with invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/chats/testAppToken", bytes.NewReader([]byte(`{invalid json`)))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid character")

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Application not found in Redis
	t.Run("application not found", func(t *testing.T) {
		// Mock Redis behavior
		mockRedis.On("KeyExists", mock.Anything, "applications_chats_count", "app#testAppToken").Return(false, nil)

		b, e := mockRedis.KeyExists(context.Background(), "applications_chats_count", "app#testAppToken")
		t.Logf("#####################%t %q", b, e)

		// Create the HTTP request with valid data
		chat := models.Chat{
			ApplicationToken: "testAppToken",
			Number:           1,
			Title:            "testTitle",
		}
		jsonData, _ := json.Marshal(chat)
		req, _ := http.NewRequest(http.MethodPost, "/chats/testAppToken", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Application not found", response["error"])

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Redis error on incrementing the field
	t.Run("failed to increment chat number", func(t *testing.T) {
		// Mock Redis behavior
		mockRedis.On("KeyExists", mock.Anything, "applications_chats_count", "app#testAppToken").Return(true, nil)
		mockRedis.On("IncrementField", mock.Anything, "applications_chats_count", "app#testAppToken").Return(int64(0), assert.AnError)

		// Create the HTTP request with valid data
		chat := models.Chat{
			ApplicationToken: "testAppToken",
			Title:            "TestTitle",
		}
		jsonData, _ := json.Marshal(chat)
		req, _ := http.NewRequest(http.MethodPost, "/chats/testAppToken", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to increment chat number", response["error"])

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})

	// Test case: Redis error on pushing to the queue
	t.Run("failed to push chat to Redis queue", func(t *testing.T) {
		// Mock Redis behavior
		mockRedis.On("KeyExists", mock.Anything, "applications_chats_count", "app#testAppToken").Return(true, nil)
		mockRedis.On("IncrementField", mock.Anything, "applications_chats_count", "app#testAppToken").Return(int64(1), nil)
		mockRedis.On("PushToQueue", mock.Anything, "chats_queue", mock.Anything).Return(assert.AnError)

		// Create the HTTP request with valid data
		chat := models.Chat{
			ApplicationToken: "testAppToken",
			Title:            "TestTitle",
		}
		jsonData, _ := json.Marshal(chat)
		req, _ := http.NewRequest(http.MethodPost, "/chats/testAppToken", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to push chat to Redis queue", response["error"])

		t.Cleanup(func() {
			mockRedis.ExpectedCalls = nil
		})
	})
}
