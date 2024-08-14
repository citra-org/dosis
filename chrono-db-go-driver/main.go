package main

import (
	"os"
	"fmt"
	// "strings"
	"net/http"

	"github.com/citra-org/chrono-db-go-driver/client"
	"github.com/gin-gonic/gin"
)

var dbClient *client.Client
var dbName string

func main() {
	admin := os.Getenv("ADMIN_USER")
	password := os.Getenv("ADMIN_PASSWORD")

	if admin == "" || password == "" {
		fmt.Println("Environment variables ADMIN_USER and ADMIN_PASSWORD must be set")
		return
	}

	uri := fmt.Sprintf("chrono://%s:%s@127.0.0.1:3141/test1", admin, password)
	var err error

	dbClient, dbName, err = client.Connect(uri)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer dbClient.Close()

	r := gin.Default()
	r.GET("/w", handleWrite)
	r.GET("/r/:stream", handleRead)
	r.GET("/cs/:stream", handleCreateStream)

	fmt.Println("Server listening on port 3000")
	err = r.Run(":3000")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleCreateStream(c *gin.Context) {
	stream := c.Param("stream")
	err := dbClient.CreateStream(dbName, stream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating record: %s", err)})
		return
	}
	c.String(http.StatusOK, "Create operation successful")
}


type Event struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

func handleWrite(c *gin.Context) {
	// stream := c.Param("stream")

	// var events []Event
	// if err := c.ShouldBindJSON(&events); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error decoding request body: %s", err)})
	// 	return
	// }

	// var eventStrings []string
	// for _, event := range events {
	// 	eventStrings = append(eventStrings, fmt.Sprintf(`("%s", "%s")`, event.Header, event.Body))
	// }
	// formattedData := fmt.Sprintf("{%s}", strings.Join(eventStrings, ", "))
	command := `INSERT INTO stream1 VALUES ('header1', 'body1'), ('header2', 'body2'), ('header3', 'body3')`
	err := dbClient.WriteEvent(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error writing data: %s", err)})
		return
	}
	c.String(http.StatusOK, "Write operation successful")
}


func handleRead(c *gin.Context) {
	stream := c.Param("stream")
	response, err := dbClient.Read(dbName, stream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error reading data: %s", err)})
		return
	}
	c.String(http.StatusOK, response)
}
