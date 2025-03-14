package handlers

import (
	"fmt"
	"log"

	"github.com/EraldCaka/aws-cloudtrail-log-forwarder/internal/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	redis   *services.RedisService
	aws     *services.AWSService
	mongo   *services.MongoService
	webhook services.WebhookService
}

func NewHandler(redis *services.RedisService, aws *services.AWSService, mongo *services.MongoService, webhook services.WebhookService) *Handler {
	return &Handler{
		redis:   redis,
		aws:     aws,
		mongo:   mongo,
		webhook: webhook,
	}
}

func (h *Handler) GetSources(c *fiber.Ctx) error {
	sources, err := h.mongo.GetSources()
	if err != nil {
		log.Printf("❌ Error fetching sources: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error fetching sources: %v", err))
	}

	log.Println("Sources retrieved successfully.")
	return c.Status(fiber.StatusOK).JSON(sources)
}

func (h *Handler) AddSource(c *fiber.Ctx) error {
	var source services.Source
	if err := c.BodyParser(&source); err != nil {
		log.Printf("❌ Invalid source data: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Invalid source data: %v", err))
	}

	err := h.mongo.InsertSource(source)
	if err != nil {
		log.Printf("❌ Error inserting source: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error inserting source: %v", err))
	}

	log.Printf("✅ Source added successfully: %v", source)
	return c.Status(fiber.StatusCreated).SendString("Source added successfully")
}

func (h *Handler) DeleteSource(c *fiber.Ctx) error {
	sourceID := c.Params("id")

	err := h.mongo.RemoveSource(sourceID)
	if err != nil {
		log.Printf("❌ Error deleting source with ID %s: %v", sourceID, err)
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error deleting source: %v", err))
	}

	log.Printf("✅ Source with ID %s deleted successfully.", sourceID)
	return c.Status(fiber.StatusOK).SendString("Source deleted successfully")
}
