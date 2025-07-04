package message_service

import (
	"context"
	"errors"
	"github.com/hasElvin/messenger-svc/internal/core/domain"
	"github.com/hasElvin/messenger-svc/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSendMessage_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data
	message := domain.Message{ID: 1, To: "+905551111001", Content: "Test message 1"}
	expectedMessageID := "msg-12345"

	// Set up expectations
	messageSender.On("Send", ctx, message).Return(expectedMessageID, nil)
	messageRepo.On("UpdateMessageStatus", ctx, uint(1), domain.StatusSent).Return(nil)
	cacheService.On("Set", ctx, "msg:1", mock.MatchedBy(func(value string) bool {
		// Verify cache value contains messageId and sentAt
		return assert.Contains(t, value, "messageId="+expectedMessageID) &&
			assert.Contains(t, value, "sentAt=")
	})).Return(nil)

	// Act
	err := service.SendMessage(ctx, message)

	// Assert
	assert.NoError(t, err)
	messageSender.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
	cacheService.AssertExpectations(t)
}

func TestSendMessage_SendError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data
	message := domain.Message{ID: 1, To: "+905551111001", Content: "Test message 1"}
	expectedError := errors.New("send failed")

	// Set up expectations - Send will fail
	messageSender.On("Send", ctx, message).Return("", expectedError)

	// Act
	err := service.SendMessage(ctx, message)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messageSender.AssertExpectations(t)

	// UpdateMessageStatus and cache Set should not be called
	messageRepo.AssertNotCalled(t, "UpdateMessageStatus")
	cacheService.AssertNotCalled(t, "Set")
}

func TestSendMessage_UpdateStatusError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data
	message := domain.Message{ID: 1, To: "+905551111001", Content: "Test message 1"}
	expectedMessageID := "msg-12345"
	updateError := errors.New("update failed")

	// Set up expectations - UpdateMessageStatus will fail
	messageSender.On("Send", ctx, message).Return(expectedMessageID, nil)
	messageRepo.On("UpdateMessageStatus", ctx, uint(1), domain.StatusSent).Return(updateError)

	// Act
	err := service.SendMessage(ctx, message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update message status")
	messageSender.AssertExpectations(t)
	messageRepo.AssertExpectations(t)

	// Cache Set should not be called when UpdateMessageStatus fails
	cacheService.AssertNotCalled(t, "Set")
}

func TestSendMessage_CacheError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	messageRepo := new(mockedMessageRepo)
	cacheService := new(mockedCacheService)
	messageSender := new(mockedMessageSender)

	// Create service instance
	service := services.NewMessageService(messageRepo, cacheService, messageSender)

	// Mock data
	message := domain.Message{ID: 1, To: "+905551111001", Content: "Test message 1"}
	expectedMessageID := "msg-12345"
	cacheError := errors.New("cache failed")

	// Set up expectations - Caching will fail
	messageSender.On("Send", ctx, message).Return(expectedMessageID, nil)
	messageRepo.On("UpdateMessageStatus", ctx, uint(1), domain.StatusSent).Return(nil)
	cacheService.On("Set", ctx, "msg:1", mock.AnythingOfType("string")).Return(cacheError)

	// Act
	err := service.SendMessage(ctx, message)

	// Assert
	// Method should still succeed even if caching fails (cache error is only logged)
	assert.NoError(t, err)
	messageSender.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
	cacheService.AssertExpectations(t)
}
