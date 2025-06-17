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
)

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

	// Initialize service
	messageService := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Seed test data for easy testing purposes
	db.SeedSampleMessages(database)

	// Start the auto sender immediately with context
	if err := messageService.StartAutoSender(context.Background(), cfg.App.SendIntervalSecs); err != nil {
		log.Fatalf("Auto sender failed to automatically start: %v", err)
	}

	// Initialize and start HTTP server
	server := rest.NewServer(messageService)

	fmt.Println("Server starting on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
