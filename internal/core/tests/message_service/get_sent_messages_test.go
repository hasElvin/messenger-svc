package message_service

import (
	"context"
	"errors"
	"github.com/hasElvin/messenger-svc/config"
	"github.com/hasElvin/messenger-svc/internal/core/domain"
	"github.com/hasElvin/messenger-svc/internal/core/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSentMessages_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	cfg := &config.Config{}
	cfg.App.MessageCharLimit = 1000

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data
	messages := []domain.Message{
		{ID: 1, To: "+905551111001", Content: "Test message 1", Status: "sent"},
		{ID: 2, To: "+905551111001", Content: "Test message 2", Status: "sent"},
	}

	// Set up expectations
	messageRepo.On("GetSentMessages", ctx).Return(messages, nil)

	//Act
	result, err := service.GetSentMessages(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, messages, result)
	messageRepo.AssertExpectations(t)
}

func TestGetSentMessages_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	expectedError := errors.New("database error")
	messageRepo.On("GetSentMessages", ctx).Return([]domain.Message{}, expectedError)

	// Act
	result, err := service.GetSentMessages(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	messageRepo.AssertExpectations(t)
}
