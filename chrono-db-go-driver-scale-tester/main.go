package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/citra-org/chrono-db-go-driver/client"
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
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer dbClient.Close()

	fmt.Println("Connected to server. Starting load test...")

	duration := 30 * time.Second
	startTime := time.Now()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup
	var totalResponseTime time.Duration
	var totalRequests int
	requestCount := 10000

	for elapsed := time.Since(startTime); elapsed < duration; elapsed = time.Since(startTime) {
		<-ticker.C
		wg.Add(requestCount)

		for i := 1; i <= requestCount; i++ {
			go func(i int) {
				defer wg.Done()
				start := time.Now()
				sendWriteRequest(i)
				responseTime := time.Since(start)
				totalResponseTime += responseTime
				totalRequests++
			}(i)
		}
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total test duration: %s\n", elapsedTime)
	fmt.Printf("Total number of requests: %d\n", totalRequests)
	fmt.Printf("Total response time for all requests: %s\n", totalResponseTime)
}

func sendWriteRequest(requestIndex int) {
	header := fmt.Sprintf("%dHEADERBa6Mf9EfFfAgt6KzwBKHS87Lma8we0kw9VSoDHtv9pQImKne5Uf4jPBUWiOfvxxvp9i1bOTpIJfErGx3uO7mS368rVx10lkw3tyQOLKD73T7BGDeEL2zJU6NTrl53oP7Da5xkUVGx18FFbQQsIK3fI9gIdK4055pyvXV36V8c30gp97damnAs0ablYGRkyGXVlV6jiscqTbo0HqD8gYENmm2iPEVJXQrFCu8ZsXknMdLym0TA3fF7nHX5oMqLHVfvfPuzmNxJHmHNMawYFrmsTQ9qp9KY1lrc5", requestIndex)
	body := fmt.Sprintf("%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", requestIndex)
	formattedData := fmt.Sprintf(`{("%s", "%s")}`, header, body)

	err := dbClient.WriteEvent(dbName, "stream1", formattedData)
	if err != nil {
		fmt.Printf("Error writing request %d: %s\n", requestIndex, err)
	}
}
