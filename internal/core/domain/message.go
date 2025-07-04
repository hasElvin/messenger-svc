package domain

import "time"

type Message struct {
	ID         uint       `json:"id"`
	To         string     `json:"to"`
	Content    string     `json:"content"`
	Status     string     `json:"status"`
	SentAt     *time.Time `json:"sent_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	RetryCount int        `json:"retry_count"`
}

const (
	StatusPending = "pending"
	StatusSent    = "sent"
	StatusFailed  = "failed"
)
