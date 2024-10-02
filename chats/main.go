package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type user struct {
    Id        int     `json:"id"`
    Username  string  `json:"username"`
}

var users = []user{
    {Id: 546, Username: "John"},
    {Id: 894, Username: "Mary"},
    {Id: 326, Username: "Jane"},
}

func getUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}

func main() {
    router := gin.Default()

    // new `GET /users` route associated with our `getUsers` function
    router.GET("/users", getUsers)

    router.Run("0.0.0.0:8080")
}