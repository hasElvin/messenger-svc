package message_service

import (
	"context"
	"github.com/hasElvin/messenger-svc/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type mockedMessageRepo struct {
	mock.Mock
}

type mockedCacheService struct {
	mock.Mock
}

type mockedMessageSender struct {
	mock.Mock
}

func (r *mockedMessageRepo) GetPendingMessages(ctx context.Context,
	limit, messageCharLimit, maxRetries int) ([]domain.Message, error) {

	args := r.Called(ctx, limit, messageCharLimit, maxRetries)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (r *mockedMessageRepo) UpdateMessageStatus(ctx context.Context,
	id uint, status string) error {

	args := r.Called(ctx, id, status)
	return args.Error(0)
}

func (r *mockedMessageRepo) GetSentMessages(ctx context.Context) ([]domain.Message, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (r *mockedMessageRepo) CreateMessage(ctx context.Context, message *domain.Message) error {
	args := r.Called(ctx, message)
	return args.Error(0)
}

func (r *mockedMessageRepo) IncrementRetryCount(ctx context.Context, id uint) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *mockedMessageRepo) SeedSampleMessages() error {
	args := r.Called()
	return args.Error(0)
}

func (r *mockedMessageRepo) ClearDatabase() error {
	args := r.Called()
	return args.Error(0)
}

func (c *mockedCacheService) Set(ctx context.Context, key, value string) error {
	args := c.Called(ctx, key, value)
	return args.Error(0)
}

func (c *mockedCacheService) Get(ctx context.Context, key string) (string, error) {
	args := c.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (s *mockedMessageSender) Send(ctx context.Context, message domain.Message) (string, error) {
	args := s.Called(ctx, message)
	return args.String(0), args.Error(1)
}
