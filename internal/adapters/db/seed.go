package db

import (
	"gorm.io/gorm"
	"log"
)

// SeedSampleMessages inserts sample data for testing
func SeedSampleMessages(db *gorm.DB) {
	var count int64
	db.Model(&MessageModel{}).Count(&count)
	if count > 0 {
		log.Println("Sample messages already seeded")
		return
	}

	messages := []MessageModel{
		{To: "+905551111001", Content: "Test message 1", Status: "pending"},
		{To: "+905551111002", Content: "Test message 2", Status: "pending"},
		{To: "+905551111003", Content: "Test message 3", Status: "pending"},
		{To: "+905551111004", Content: "Test message 4", Status: "pending"},
		{To: "+905551111005", Content: "Test message 5", Status: "pending"},
		{To: "+905551111006", Content: "Test message 6", Status: "pending"},
		{To: "+905551111007", Content: "Test message 7", Status: "pending"},
		{To: "+905551111008", Content: "Test message 8", Status: "pending"},
		{To: "+905551111009", Content: "Test message 9", Status: "pending"},
		{To: "+905551111010", Content: "Test message 10", Status: "pending"},
	}

	if err := db.Create(&messages).Error; err != nil {
		log.Printf("Failed to seed sample messages: %v", err)
		return
	}

	log.Println("10 sample messages inserted into the database")
}
