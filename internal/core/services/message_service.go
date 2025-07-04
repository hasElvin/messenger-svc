package services

import (
	"context"
	"fmt"
	"github.com/hasElvin/messenger-svc/config"
	"log"
	"sync"
	"time"

	"github.com/hasElvin/messenger-svc/internal/core/domain"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
)

type messageService struct {
	repo      ports.MessageRepository
	cache     ports.CacheService
	sender    ports.MessageSender
	stopChan  chan struct{}
	isRunning bool
	mu        sync.RWMutex
}

func NewMessageService(repo ports.MessageRepository, cache ports.CacheService,
	sender ports.MessageSender) ports.MessageService {
	return &messageService{
		repo:   repo,
		cache:  cache,
		sender: sender,
	}
}

func (s *messageService) StartAutoSender(ctx context.Context, intervalSeconds int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("auto sender is already running")
	}

	s.stopChan = make(chan struct{})
	s.isRunning = true

	go s.runAutoSender(ctx, intervalSeconds)
	log.Println("Auto sender started")

	return nil
}

func (s *messageService) StopAutoSender(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("auto sender is not running")
	}

	close(s.stopChan)
	s.isRunning = false
	log.Println("Auto sender stopped")

	return nil
}

func (s *messageService) GetSentMessages(ctx context.Context) ([]domain.Message, error) {
	return s.repo.GetSentMessages(ctx)
}

func (s *messageService) runAutoSender(ctx context.Context, intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	cfg := config.LoadConfig()

	for {
		select {
		case <-ticker.C:
			s.SendPendingMessages(ctx, &cfg)
		case <-s.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *messageService) SendPendingMessages(ctx context.Context, cfg *config.Config) {
	messages, err := s.repo.GetPendingMessages(ctx, 2, cfg.App.MessageCharLimit, cfg.App.MaxRetries)
	if err != nil {
		log.Printf("Failed to fetch pending messages: %v", err)
		return
	}

	for _, msg := range messages {
		err := s.SendMessage(ctx, msg)
		if err != nil {
			log.Printf("Failed to send message ID %d: %v", msg.ID, err)

			msg.RetryCount++
			if msg.RetryCount >= cfg.App.MaxRetries {
				log.Printf("Marking message ID %d as failed after %d retries", msg.ID, msg.RetryCount)
				_ = s.repo.IncrementRetryCount(ctx, msg.ID)
				_ = s.repo.UpdateMessageStatus(ctx, msg.ID, domain.StatusFailed)
			} else {
				_ = s.repo.IncrementRetryCount(ctx, msg.ID)
			}
		}
	}
}

func (s *messageService) SendMessage(ctx context.Context, msg domain.Message) error {
	messageID, err := s.sender.Send(ctx, msg)
	if err != nil {
		return err
	}

	// Update message status
	if err := s.repo.UpdateMessageStatus(ctx, msg.ID, domain.StatusSent); err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	// Cache the result
	cacheKey := fmt.Sprintf("msg:%d", msg.ID)
	cacheValue := fmt.Sprintf("messageId=%s|sentAt=%s", messageID, time.Now().Format(time.RFC3339))

	if err := s.cache.Set(ctx, cacheKey, cacheValue); err != nil {
		log.Printf("Failed to cache message %d: %v", msg.ID, err)
	}

	log.Printf("Message %d sent successfully", msg.ID)
	return nil
}
