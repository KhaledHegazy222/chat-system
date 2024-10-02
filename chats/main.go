package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "fmt"
	"net/http"
	"os"
	"strconv"

	// "strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

// Database connection pool
var db *sql.DB
var redis_client *redis.Client

func getChats(c *gin.Context) {
	rows, err := db.Query("SELECT chats.number, applications.token FROM applications join chats on chats.application_id = applications.id")

	if err != nil {
		print(err.Error())
		return
	}
	defer rows.Close()

	chats := []chat{}

	// Iterate through the result set and scan into chat structs
	for rows.Next() {
		var chat chat
		if err := rows.Scan(&chat.Number, &chat.ApplicationToken); err != nil {
			print(err.Error())
			return
		}
		chats = append(chats, chat)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.IndentedJSON(http.StatusOK, chats)
}

func createChat(c *gin.Context) {
	var newChat chat

	// Parse the JSON request body into the chat struct
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	err = redis_client.RPush(ctx, "test_chats_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push chat to Redis queue"})
		return
	}

	// Return the created chat data immediately
	c.JSON(http.StatusCreated, newChat)
}
func updateChat(c *gin.Context) {
	// Get the application_token and chat_number from the URL parameters
	applicationToken := c.Param("application_token")
	chatNumber := c.Param("chat_number")

	// Parse the JSON request body into the chat struct
	var updatedChat chat
	if err := c.ShouldBindJSON(&updatedChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the application_token and chat_number from the URL params
	updatedChat.ApplicationToken = applicationToken

	// Convert chat_number to int
	chatNum, err := strconv.Atoi(chatNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat number"})
		return
	}
	updatedChat.Number = chatNum

	// Create a queue message with the operation type "update"
	queueData := queueMessage{
		Operation: "update", // Set the operation type to update
		Data:      updatedChat,
	}

	// Serialize the queue message to JSON
	queueMessageData, err := json.Marshal(queueData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize chat data"})
		return
	}

	// Push the serialized chat data onto a Redis list (queue)
	ctx := context.Background()
	err = redis_client.RPush(ctx, "test_chats_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push chat to Redis queue"})
		return
	}

	// Return the updated chat data immediately
	c.JSON(http.StatusOK, updatedChat)
}

func deleteChat(c *gin.Context) {
	// Get the application_token and chat_number from the URL parameters
	applicationToken := c.Param("application_token")
	chatNumber := c.Param("chat_number")

	// Create a chat struct using the provided parameters
	var chatToDelete chat
	chatToDelete.ApplicationToken = applicationToken

	// Convert chat_number to int
	chatNum, err := strconv.Atoi(chatNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat number"})
		return
	}
	chatToDelete.Number = chatNum

	// Create a queue message with the operation type "delete"
	queueData := queueMessage{
		Operation: "delete", // Set the operation type to delete
		Data:      chatToDelete,
	}

	// Serialize the queue message to JSON
	queueMessageData, err := json.Marshal(queueData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize chat data"})
		return
	}

	// Push the serialized chat data onto a Redis list (queue)
	ctx := context.Background()
	err = redis_client.RPush(ctx, "test_chats_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push chat to Redis queue"})
		return
	}

	// Return a confirmation message
	c.JSON(http.StatusOK, gin.H{"status": "Chat delete request added to queue"})
}

func getMessages(c *gin.Context) {
	// Prepare SQL query
	rows, err := db.Query("SELECT messages.number, chats.number, applications.token,messages.content FROM applications join chats on chats.application_id = applications.id join messages on chats.id = messages.chat_id")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	messages := []message{}

	// Iterate through the result set and scan into message structs
	for rows.Next() {
		var msg message
		if err := rows.Scan(&msg.Number, &msg.ChatNumber, &msg.ApplicationToken, &msg.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messages = append(messages, msg)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.IndentedJSON(http.StatusOK, messages)
}

func createMessage(c *gin.Context) {
	var newMessage message

	// Parse the JSON request body into the message struct
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	err = redis_client.RPush(ctx, "test_messages_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push message to Redis queue"})
		return
	}

	// Return the created message data immediately
	c.JSON(http.StatusCreated, newMessage)
}

func updateMessage(c *gin.Context) {
	// Get the application_token and chat_number from the URL parameters
	applicationToken := c.Param("application_token")
	chatNumber := c.Param("chat_number")
	messageNumber := c.Param("message_number")

	// Parse the JSON request body into the message struct
	var updatedMessage message
	if err := c.ShouldBindJSON(&updatedMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the application_token and chat_number from the URL params
	updatedMessage.ApplicationToken = applicationToken

	// Convert chat_number to int
	chatNum, err := strconv.Atoi(chatNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat number"})
		return
	}
	updatedMessage.ChatNumber = chatNum

	// Convert chat_number to int
	messageNum, err := strconv.Atoi(messageNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat number"})
		return
	}
	updatedMessage.Number = messageNum

	// Create a queue message with the operation type "update"
	queueData := queueMessage{
		Operation: "update", // Set the operation type to update
		Data:      updatedMessage,
	}

	// Serialize the queue message to JSON
	queueMessageData, err := json.Marshal(queueData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize message data"})
		return
	}

	// Push the serialized message data onto a Redis list (queue)
	ctx := context.Background()
	err = redis_client.RPush(ctx, "test_messages_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push message to Redis queue"})
		return
	}

	// Return the updated message data immediately
	c.JSON(http.StatusOK, updatedMessage)
}
func deleteMessage(c *gin.Context) {
	// Get the application_token, chat_number, and message_number from the URL parameters
	applicationToken := c.Param("application_token")
	chatNumber := c.Param("chat_number")
	messageNumber := c.Param("message_number")

	// Create a message struct using the provided parameters
	var messageToDelete message
	messageToDelete.ApplicationToken = applicationToken

	// Convert chat_number to int
	chatNum, err := strconv.Atoi(chatNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat number"})
		return
	}
	messageToDelete.ChatNumber = chatNum

	// Convert message_number to int
	messageNum, err := strconv.Atoi(messageNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message number"})
		return
	}
	messageToDelete.Number = messageNum

	// Create a queue message with the operation type "delete"
	queueData := queueMessage{
		Operation: "delete", // Set the operation type to delete
		Data:      messageToDelete,
	}

	// Serialize the queue message to JSON
	queueMessageData, err := json.Marshal(queueData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize message data"})
		return
	}

	// Push the serialized message data onto a Redis list (queue)
	ctx := context.Background()
	err = redis_client.RPush(ctx, "test_messages_queue", queueMessageData).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push message to Redis queue"})
		return
	}

	// Return a confirmation message
	c.JSON(http.StatusOK, gin.H{"status": "Message delete request added to queue"})
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

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	var err error
	// Connect Database Server
	db, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := gin.Default()

	// Chat routes
	router.GET("/chats", getChats)
	router.POST("/chats", createChat)
	router.PATCH("/chats/:application_token/:chat_number", updateChat)
	router.DELETE("/chats/:application_token/:chat_number", deleteChat)

	// Message routes
	router.GET("/messages", getMessages)
	router.POST("/messages", createMessage)
	router.PATCH("/messages/:application_token/:chat_number/:message_number", updateMessage)
	router.DELETE("/messages/:application_token/:chat_number/:message_number", deleteMessage)

	router.Run("0.0.0.0:8080")

}
