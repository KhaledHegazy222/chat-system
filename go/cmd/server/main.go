package main

import (
	"github.com/gin-gonic/gin"
	"chat_system/internal/api"
	"chat_system/internal/redis"
)

func main() {
	// Initialize Redis
	redis.InitRedis()

	// Set up the Gin router
	router := gin.Default()

	// Chat routes
	router.POST("/applications/:application_token/chats", api.CreateChat)

	// Message routes
	router.POST("/applications/:application_token/chats/:chat_number/messages", api.CreateMessage)

	// Start the server
	router.Run("0.0.0.0:8080")
}
