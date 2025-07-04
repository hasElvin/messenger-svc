package message_service

import (
	"context"
	"errors"
	"github.com/hasElvin/messenger-svc/config"
	"github.com/hasElvin/messenger-svc/internal/core/domain"
	"github.com/hasElvin/messenger-svc/internal/core/services"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSendPendingMessages_Success(t *testing.T) {
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
		{ID: 1, To: "+905551111001", Content: "Test message 1"},
		{ID: 2, To: "+905551111001", Content: "Test message 2"},
	}

	// Set up expectations
	messageRepo.On("GetPendingMessages", ctx, 2, mock.Anything, mock.Anything).Return(messages, nil)

	// Mock sendMessage calls (these will be called for each message)
	messageSender.On("Send", ctx, messages[0]).Return("msg-id-1", nil)
	messageSender.On("Send", ctx, messages[1]).Return("msg-id-2", nil)

	messageRepo.On("UpdateMessageStatus", ctx, uint(1), domain.StatusSent).Return(nil)
	messageRepo.On("UpdateMessageStatus", ctx, uint(2), domain.StatusSent).Return(nil)

	cacheService.On("Set", ctx, "msg:1", mock.AnythingOfType("string")).Return(nil)
	cacheService.On("Set", ctx, "msg:2", mock.AnythingOfType("string")).Return(nil)

	// Act
	service.SendPendingMessages(ctx, cfg)

	// Assert
	messageRepo.AssertExpectations(t)
	messageSender.AssertExpectations(t)
	cacheService.AssertExpectations(t)
}

func TestSendPendingMessages_GetPendingMessagesError(t *testing.T) {
	// Arrange
	ctx := context.Background()

	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	cfg := &config.Config{}
	cfg.App.MessageCharLimit = 1000

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Set up expectations - repo returns error
	messageRepo.On("GetPendingMessages", ctx, 2, mock.Anything, mock.Anything).
		Return([]domain.Message{}, errors.New("database error"))

	// Act
	service.SendPendingMessages(ctx, cfg)

	// Assert
	messageRepo.AssertExpectations(t)

	// sendMessage should not be called when GetPendingMessages fails
	messageSender.AssertNotCalled(t, "Send")
}

func TestSendPendingMessages_SendMessageError_SingleMessage(t *testing.T) {
	// Arrange
	ctx := context.Background()

	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	cfg := &config.Config{}
	cfg.App.MessageCharLimit = 1000
	cfg.App.MaxRetries = 3

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data - message with RetryCount = 0 (first attempt)
	messages := []domain.Message{
		{ID: 1, To: "+905551111001", Content: "Test message 1", RetryCount: 0},
	}

	// Set up expectations - sendMessage will fail
	messageRepo.On("GetPendingMessages", ctx, 2, mock.Anything, mock.Anything).Return(messages, nil)
	messageSender.On("Send", ctx, messages[0]).Return("", errors.New("send error"))

	messageRepo.On("IncrementRetryCount", ctx, uint(1)).Return(nil)

	// Act
	service.SendPendingMessages(ctx, cfg)

	// Assert
	messageRepo.AssertExpectations(t)
	messageSender.AssertExpectations(t)

	// UpdateMessageStatus with StatusFailed should not be called (not max retries yet)
	messageRepo.AssertNotCalled(t, "UpdateMessageStatus")

	// Cache Set should not be called when Send fails
	cacheService.AssertNotCalled(t, "Set")
}

func TestSendPendingMessages_SendMessageError_MultipleMessages(t *testing.T) {
	// Arrange
	ctx := context.Background()

	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	cfg := &config.Config{}
	cfg.App.MessageCharLimit = 1000
	cfg.App.MaxRetries = 3

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data
	messages := []domain.Message{
		{ID: 1, To: "+905551111001", Content: "Test message 1", RetryCount: 0}, // This will fail
		{ID: 2, To: "+905551111002", Content: "Test message 2", RetryCount: 0}, // This will succeed
	}

	// Set up expectations - sendMessage will fail for first message
	messageRepo.On("GetPendingMessages", ctx, 2, mock.Anything, mock.Anything).Return(messages, nil)
	messageSender.On("Send", ctx, messages[0]).Return("", errors.New("send error"))
	messageSender.On("Send", ctx, messages[1]).Return("msg-id-2", nil)

	// First message fails - only IncrementRetryCount should be called (not at max retries)
	messageRepo.On("IncrementRetryCount", ctx, uint(1)).Return(nil)

	// The second message should still update status and cache
	messageRepo.On("UpdateMessageStatus", ctx, uint(2), domain.StatusSent).Return(nil)
	cacheService.On("Set", ctx, "msg:2", mock.AnythingOfType("string")).Return(nil)

	// Act
	service.SendPendingMessages(ctx, cfg)

	// Assert
	messageRepo.AssertExpectations(t)
	messageSender.AssertExpectations(t)
	cacheService.AssertExpectations(t)

	// UpdateMessageStatus and cache Set should not be called for the first message
	messageRepo.AssertNotCalled(t, "UpdateMessageStatus", ctx, uint(1), domain.StatusSent)
	cacheService.AssertNotCalled(t, "Set", ctx, "msg:1", mock.AnythingOfType("string"))
}

func TestSendPendingMessages_SendMessageError_MaxRetriesReached(t *testing.T) {
	// Arrange
	ctx := context.Background()

	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	cfg := &config.Config{}
	cfg.App.MessageCharLimit = 1000
	cfg.App.MaxRetries = 3

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data - message with RetryCount = 2 (will reach max retries after increment)
	messages := []domain.Message{
		{ID: 1, To: "+905551111001", Content: "Test message 1", RetryCount: 2},
	}

	// Set up expectations - sendMessage will fail
	messageRepo.On("GetPendingMessages", ctx, 2, mock.Anything, mock.Anything).Return(messages, nil)
	messageSender.On("Send", ctx, messages[0]).Return("", errors.New("send error"))

	messageRepo.On("IncrementRetryCount", ctx, uint(1)).Return(nil)
	messageRepo.On("UpdateMessageStatus", ctx, uint(1), domain.StatusFailed).Return(nil)

	// Act
	service.SendPendingMessages(ctx, cfg)

	// Assert
	messageRepo.AssertExpectations(t)
	messageSender.AssertExpectations(t)

	// Cache Set should not be called when Send fails
	cacheService.AssertNotCalled(t, "Set")
}
