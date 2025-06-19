package main

import (
	"context"
	"fmt"
	"github.com/hasElvin/messenger-svc/config"
	"github.com/hasElvin/messenger-svc/internal/adapters/cache"
	"github.com/hasElvin/messenger-svc/internal/adapters/db"
	"github.com/hasElvin/messenger-svc/internal/adapters/http"
	"github.com/hasElvin/messenger-svc/internal/adapters/rest"
	"github.com/hasElvin/messenger-svc/internal/core/services"
	"log"
	"os"
)

// @title Messenger API
// @version 1.0
// @description API for managing message auto-sending system
// @host messenger-svc-gfsy.onrender.com
// @schemes https
// @BasePath /
func main() {
	cfg := config.LoadConfig()

	// Initialize database
	database := db.InitPostgres(cfg)

	// Initialize Redis
	redisClient := cache.InitRedis(cfg)

	// Initialize adapters
	messageRepo := db.NewPostgresRepository(database)
	cacheService := cache.NewRedisCache(redisClient)
	messageSender := http.NewWebhookSender(cfg.App.WebhookURL)

	// Initialize services
	messageService := services.NewMessageService(messageRepo, cacheService, messageSender)
	utilityService := services.NewUtilityService(messageRepo)

	// Seed test data for easy testing purposes
	if err := utilityService.SeedSampleMessages(); err != nil {
		log.Fatalf("Failed to seed sample messages: %v", err)
	}

	// Start the auto sender immediately with context
	if err := messageService.StartAutoSender(context.Background(), cfg.App.SendIntervalSecs); err != nil {
		log.Fatalf("Auto sender failed to automatically start: %v", err)
	}

	// Initialize and start HTTP server
	server := rest.NewServer(messageService, utilityService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on :%s\n", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
