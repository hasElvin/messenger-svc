package ports

import (
	"context"
	"github.com/hasElvin/messenger-svc/internal/core/domain"
)

// MessageRepository defines the interface for message persistence
type MessageRepository interface {
	GetPendingMessages(ctx context.Context, limit int) ([]domain.Message, error)
	UpdateMessageStatus(ctx context.Context, id uint, status string) error
	GetSentMessages(ctx context.Context) ([]domain.Message, error)
	CreateMessage(ctx context.Context, message *domain.Message) error
	SeedSampleMessages() error
	ClearDatabase() error
}

// MessageSender defines the interface for sending messages
type MessageSender interface {
	Send(ctx context.Context, message domain.Message) (string, error)
}

// MessageService defines the interface for message business logic
type MessageService interface {
	StartAutoSender(ctx context.Context, intervalSeconds int) error
	StopAutoSender(ctx context.Context) error
	GetSentMessages(ctx context.Context) ([]domain.Message, error)
}

// UtilityService defines some utility tools for testing the app
type UtilityService interface {
	SeedSampleMessages() error
	ClearDatabase() error
}
