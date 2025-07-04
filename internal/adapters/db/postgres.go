package db

import (
	"fmt"
	"github.com/hasElvin/messenger-svc/config"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type postgresRepository struct {
	db *gorm.DB
}

func InitPostgres(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name, cfg.Database.SSLMode,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the message table
	if err := database.AutoMigrate(&MessageModel{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return database
}

func NewPostgresRepository(db *gorm.DB) ports.MessageRepository {
	return &postgresRepository{db: db}
}
