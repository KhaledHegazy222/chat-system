package api

import (
	"chat_system/internal/models"
	"chat_system/internal/redis"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateChat(c *gin.Context, redis redis.RedisClient) {
	var newChat models.Chat
	appToken := c.Param("application_token")
	newChat.ApplicationToken = appToken

	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation for the required Title field
	if newChat.Title == "" { // Assuming Title is the only required field
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is a required field"})
		return
	}

	// Trim whitespace from the Title field
	newChat.Title = strings.TrimSpace(newChat.Title)

	ctx := c.Request.Context()
	appHashKey := "app#" + appToken

	exists, err := redis.KeyExists(ctx, "applications_chats_count", appHashKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check application existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	chatNumber, err := redis.IncrementField(ctx, "applications_chats_count", appHashKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment chat number"})
		return
	}
	newChat.Number = int(chatNumber)

	if err := redis.PushToQueue(ctx, "chats_queue", newChat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push chat to Redis queue"})
		return
	}

	c.JSON(http.StatusCreated, newChat)
}
