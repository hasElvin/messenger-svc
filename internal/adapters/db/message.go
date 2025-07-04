package db

import (
	"context"
	"github.com/hasElvin/messenger-svc/internal/core/domain"
	"gorm.io/gorm"
	"time"
)

type MessageModel struct {
	ID         uint   `gorm:"primaryKey"`
	To         string `gorm:"not null"`
	Content    string `gorm:"not null;size:160"`
	Status     string `gorm:"default:'pending'"`
	RetryCount int    `gorm:"default:0"`
	SentAt     *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (MessageModel) TableName() string {
	return "messages"
}

func (r *postgresRepository) GetPendingMessages(ctx context.Context,
	limit, messageCharLimit, maxRetries int) ([]domain.Message, error) {

	var models []MessageModel
	err := r.db.WithContext(ctx).
		Where("status = ? AND char_length(content) <= ? AND retry_count < ?",
			domain.StatusPending, messageCharLimit, maxRetries).
		Limit(limit).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	messages := make([]domain.Message, 0)
	for _, model := range models {
		messages = append(messages, r.toDomain(model))
	}

	return messages, nil
}

func (r *postgresRepository) UpdateMessageStatus(ctx context.Context, id uint, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}

	if status == domain.StatusSent {
		updates["sent_at"] = now
	}

	return r.db.WithContext(ctx).
		Model(&MessageModel{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *postgresRepository) GetSentMessages(ctx context.Context) ([]domain.Message, error) {
	var models []MessageModel
	err := r.db.WithContext(ctx).
		Where("status = ?", domain.StatusSent).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	messages := make([]domain.Message, len(models))
	for i, model := range models {
		messages[i] = r.toDomain(model)
	}

	return messages, nil
}

func (r *postgresRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	model := r.toModel(*message)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return err
	}

	message.ID = model.ID
	message.CreatedAt = model.CreatedAt
	message.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *postgresRepository) IncrementRetryCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&MessageModel{}).
		Where("id = ?", id).
		Update("retry_count", gorm.Expr("retry_count + 1")).Error
}

func (r *postgresRepository) toDomain(model MessageModel) domain.Message {
	return domain.Message{
		ID:         model.ID,
		To:         model.To,
		Content:    model.Content,
		Status:     model.Status,
		SentAt:     model.SentAt,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
		RetryCount: model.RetryCount,
	}
}

func (r *postgresRepository) toModel(message domain.Message) MessageModel {
	return MessageModel{
		ID:         message.ID,
		To:         message.To,
		Content:    message.Content,
		Status:     message.Status,
		SentAt:     message.SentAt,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
		RetryCount: message.RetryCount,
	}
}
