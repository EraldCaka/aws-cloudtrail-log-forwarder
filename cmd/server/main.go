package main

import (
	"context"
	"log"
	"time"

	"github.com/EraldCaka/aws-cloudtrail-log-forwarder/internal/handlers"
	"github.com/EraldCaka/aws-cloudtrail-log-forwarder/internal/services"
	tools "github.com/EraldCaka/aws-cloudtrail-log-forwarder/internal/util"
	"github.com/gofiber/fiber/v2"
)

func main() {
	envData := services.InitEnvData()
	redisService := services.NewRedisService()
	awsService, err := services.NewAWSService(envData.AccessKey, envData.SecretKey)
	if err != nil {
		log.Fatalf("❌ Failed to initialize AWS service: %v", err)
	}
	mongoService, err := services.NewMongoService(envData.MongoDbCon, "cloudtrail_logs")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	webhookService := services.NewWebhookService()

	handler := handlers.NewHandler(redisService, awsService, mongoService, webhookService)

	app := fiber.New()

	app.Get("/sources", handler.GetSources)
	app.Post("/sources", handler.AddSource)
	app.Delete("/sources/:id", handler.DeleteSource)

	go realTimeLogCheckAndDelivery(awsService, webhookService)

	if err := app.Listen(":5005"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
func realTimeLogCheckAndDelivery(awsService *services.AWSService, webhookService services.WebhookService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			startTime := time.Now().Add(-1 * time.Minute).Format(time.RFC3339)
			endTime := time.Now().Format(time.RFC3339)

			logs, err := awsService.FetchLogs(context.Background(), startTime, endTime)
			if err != nil {
				log.Printf("Failed to fetch logs: %v", err)
				continue
			}
			for _, logEvent := range logs {
				err := webhookService.SendLog("http://localhost:8080/webhook", tools.ConvertEventToMap(logEvent))
				if err != nil {
					log.Printf("Failed to send log to webhook: %v", err)
					continue
				}
				log.Printf("Log forwarded successfully: %s", *logEvent.EventId)
			}
		}
	}
}
