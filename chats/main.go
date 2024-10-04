package main

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type chat struct {
	Number           int    `json:"number"`
	ApplicationToken string `json:"application_token"`
}

type message struct {
	Number           int    `json:"number"`
	ChatNumber       int    `json:"chat_number"`
	ApplicationToken string `json:"application_token"`
	Content          string `json:"content"`
}

type queueMessage struct {
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

// Redis Client
var redis_client *redis.Client

func createChat(c *gin.Context) {
	var newChat chat

	// Parse the JSON request body into the chat struct
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application_token := newChat.ApplicationToken
	application_name_in_hash := "app#" + application_token

	// Check if the application is already created
	// if created it will have a the application token will exists in applications_chats_count hash
	exists, err := redis_client.HExists(c, "applications_chats_count", application_name_in_hash).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existence for application"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application is not found"})
		return
	}

	// Atomically Increment The Chats Number Per Application and get the total count (chat number)
	chats_number_result, err := redis_client.HIncrBy(c, "applications_chats_count", application_name_in_hash, 1).Result()
	if err != nil {
		panic(err)
	}

	newChat.Number = int(chats_number_result)

	// Create a queue message with the operation type "create"
	queueData := queueMessage{
		Operation: "create",
		Data:      newChat,
	}

	// Serialize the queue message to JSON
	queueMessageData, err := json.Marshal(queueData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize chat data"})
		return
	}

	// Push the serialized chat data onto a Redis list (queue)
	ctx := context.Background()
	err = redis_client.RPush(ctx, "chats_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push chat to Redis queue"})
		return
	}

	// Return the created chat data immediately
	c.JSON(http.StatusCreated, newChat)
}

func createMessage(c *gin.Context) {
	var newMessage message

	// Parse the JSON request body into the message struct
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application_token := newMessage.ApplicationToken
	chat_number := newMessage.ChatNumber
	chat_name_in_hash := "chat#" + application_token + "-" + strconv.Itoa(chat_number)

	// Check if the chat is already created
	// if created it will have a the chat_name will exists in chats_messages_count hash
	exists, err := redis_client.HExists(c, "chats_messages_count", chat_name_in_hash).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existence for chat"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat is not found"})
		return
	}

	// Atomically Increment The Messages Number Per Chat and get the total count (message number)
	messages_number_result, err := redis_client.HIncrBy(c, "chats_messages_count", chat_name_in_hash, 1).Result()
	if err != nil {
		panic(err)
	}

	newMessage.Number = int(messages_number_result)

	// Create a queue message with the operation type "create"
	queueData := queueMessage{
		Operation: "create",
		Data:      newMessage,
	}

	// Serialize the queue message to JSON
	queueMessageData, err := json.Marshal(queueData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize message data"})
		return
	}

	// Push the serialized message data onto a Redis list (queue)
	ctx := context.Background()
	err = redis_client.RPush(ctx, "messages_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push message to Redis queue"})
		return
	}

	// Return the created message data immediately
	c.JSON(http.StatusCreated, newMessage)
}


func main() {
	// Connect to Redis Server
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	redis_client = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	router := gin.Default()

	// Chat routes
	router.POST("/chats", createChat)

	// Message routes
	router.POST("/messages/", createMessage)
	
	// router.PATCH("/messages/:application_token/:chat_number/:message_number", updateMessage)

	router.Run("0.0.0.0:8080")

}
