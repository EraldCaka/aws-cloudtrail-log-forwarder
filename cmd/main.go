package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EraldCaka/aws-cloudtrail-log-forwarder/internal/services"
)

func main() {
	redisService := services.NewRedisService()
	fmt.Println("✅ Redis service initialized successfully!")

	err := redisService.SetKey("test_key", "test_value")
	if err != nil {
		log.Fatalf("❌ Failed to set key in Redis: %v", err)
	}
	fmt.Println("✅ Successfully set key in Redis!")

	val, err := redisService.GetKey("test_key")
	if err != nil {
		log.Fatalf("❌ Failed to get key from Redis: %v", err)
	}
	fmt.Printf("✅ Retrieved key from Redis: %s\n", val)

	err = redisService.EnqueueTask("queue", val)
	if err != nil {
		log.Fatalf("❌ Failed to enqueue task: %v", err)
	}
	fmt.Println("✅ Task enqueued successfully in Redis!")

	err = redisService.DeleteKey("test_key")
	if err != nil {
		log.Fatalf("❌ Failed to delete key from Redis: %v", err)
	}
	fmt.Println("✅ Successfully deleted key from Redis!")

	mongoService, err := services.NewMongoService("mongodb://root:rootpassword@localhost:27017", "cloudtrail_logs")
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	fmt.Println("✅ MongoDB connected successfully!")

	testDocument := map[string]interface{}{
		"event": "AWS Login",
		"user":  "admin",
		"time":  time.Now(),
	}
	insertedID, err := mongoService.InsertDocument("logs", testDocument)
	if err != nil {
		log.Fatalf("❌ Failed to insert document in MongoDB: %v", err)
	}
	fmt.Printf("✅ Successfully inserted document in MongoDB! ID: %v\n", insertedID)

	awsService, err := services.NewAWSService()
	if err != nil {
		log.Fatalf("❌ Failed to initialize AWS service: %v", err)
	}
	fmt.Println("✅ AWS service initialized successfully!")

	startTime := "2024-03-12T14:00:00Z"
	endTime := "2024-03-12T15:00:00Z"
	events, err := awsService.FetchLogs(context.Background(), startTime, endTime)
	if err != nil {
		log.Fatalf("❌ Failed to fetch AWS CloudTrail logs: %v", err)
	}
	fmt.Printf("✅ Retrieved %d events from CloudTrail!\n", len(events))

	for _, event := range events {
		fmt.Printf("🔹 Event ID: %s | Event Name: %s | Event Time: %s\n", *event.EventId, *event.EventName, *event.EventTime)
	}

	fmt.Println("🚀 All tests completed successfully!")
}
