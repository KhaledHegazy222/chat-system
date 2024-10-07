package main

import (
	"chat_system/internal/api"
	"chat_system/internal/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Redis
	redisClient := redis.NewRedisClientWrapper()

	// Set up the Gin router
	router := gin.Default()

	// Chat routes
	router.POST("/applications/:application_token/chats", func(c *gin.Context) {
		api.CreateChat(c, redisClient)
	})

	// Message routes
	router.POST("/applications/:application_token/chats/:chat_number/messages", func(c *gin.Context) {
		api.CreateMessage(c, redisClient)
	})

	// Start the server
	router.Run("0.0.0.0:8080")
}
