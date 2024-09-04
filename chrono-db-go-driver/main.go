package main

import (
	// "os"
	"fmt"
	"time"
	// "strings"
	"net/http"

	"github.com/citra-org/chrono-db-go-driver/client"
	"github.com/gin-gonic/gin"
)

var dbClient *client.Client
var dbName string
var err error

func connectToDatabase(uri string) error {
	dbClient, dbName, err = client.Connect(uri)
	if err != nil {
		return err
	}
	return nil
}

func ensureConnection(uri string) error {
	err := dbClient.PingChrono()
	if err != nil {
		fmt.Println("Connection lost. Attempting to reconnect...")
		for {
			err := connectToDatabase(uri)
			if err != nil {
				fmt.Println("Reconnection failed:", err)
				fmt.Println("Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Println("Reconnected to the database.")
			break
		}
	}
	return nil
}

func main() {
	admin := "admin"
	password := "X%CQXZpqWOmomvIp"
	if admin == "" || password == "" {
		fmt.Println("Environment variables ADMIN_USER and ADMIN_PASSWORD must be set")
		return
	}

	uri := fmt.Sprintf("chrono://%s:%s@127.0.0.1:3141/test1", admin, password)

	for {
		err := connectToDatabase(uri)
		if err != nil {
			fmt.Println("Error connecting to database:", err)
			fmt.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	defer dbClient.Close()

	r := gin.Default()

	r.GET("/w", func(c *gin.Context) {
		if err := ensureConnection(uri); err != nil {
			c.JSON(500, gin.H{"error": "Database connection failed"})
			return
		}
		handleWrite(c)
	})

	r.GET("/r/:stream", func(c *gin.Context) {
		if err := ensureConnection(uri); err != nil {
			c.JSON(500, gin.H{"error": "Database connection failed"})
			return
		}
		handleRead(c)
	})

	r.GET("/cs/:stream", func(c *gin.Context) {
		if err := ensureConnection(uri); err != nil {
			c.JSON(500, gin.H{"error": "Database connection failed"})
			return
		}
		handleCreateStream(c)
	})

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
