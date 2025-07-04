package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hasElvin/messenger-svc/internal/core/domain"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
)

type webhookSender struct {
	webhookURL string
	client     *http.Client
}

func NewWebhookSender(webhookURL string) ports.MessageSender {
	return &webhookSender{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (w *webhookSender) Send(ctx context.Context, message domain.Message) (string, error) {
	payload := map[string]interface{}{
		"to":      message.To,
		"content": message.Content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	messageID, ok := response["messageId"].(string)
	if !ok {
		return "", fmt.Errorf("messageId not found in response")
	}

	return messageID, nil
}
