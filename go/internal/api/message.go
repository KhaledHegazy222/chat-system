package api

import (
	"chat_system/internal/models"
	"chat_system/internal/redis"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateMessage(c *gin.Context, redis redis.RedisClient) {
	var newMessage models.Message
	appToken := c.Param("application_token")
	chatNumber, err := strconv.Atoi(c.Param("chat_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat number"})
		return
	}
	newMessage.ApplicationToken = appToken
	newMessage.ChatNumber = chatNumber

	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	chatHashKey := "chat#" + appToken + "-" + strconv.Itoa(chatNumber)

	exists, err := redis.KeyExists(ctx, "chats_messages_count", chatHashKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chat existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	messageNumber, err := redis.IncrementField(ctx, "chats_messages_count", chatHashKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment message number"})
		return
	}
	newMessage.Number = int(messageNumber)

	if err := redis.PushToQueue(ctx, "messages_queue", newMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push message to Redis queue"})
		return
	}

	c.JSON(http.StatusCreated, newMessage)
}
