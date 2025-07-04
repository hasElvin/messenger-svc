package db

import (
	"log"
)

// SeedSampleMessages inserts sample data for testing
func (r *postgresRepository) SeedSampleMessages() error {
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

	if err := r.db.Create(&messages).Error; err != nil {
		log.Printf("Failed to seed sample messages: %v", err)
		return err
	}

	log.Println("10 sample messages inserted into the database")
	return nil
}

// ClearDatabase truncates the messages table to clear all data
func (r *postgresRepository) ClearDatabase() error {
	if err := r.db.Exec("TRUNCATE TABLE messages").Error; err != nil {
		log.Printf("Failed to clear database: %v", err)
		return err
	}
	log.Println("Database cleared successfully")
	return nil
}
