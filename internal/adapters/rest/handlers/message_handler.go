package handlers

import (
	"context"
	"github.com/hasElvin/messenger-svc/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
)

type MessageHandler struct {
	messageService ports.MessageService
}

func NewMessageHandler(messageService ports.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *MessageHandler) StartAutoSender(c *gin.Context) {
	cfg := config.LoadConfig()
	intervalSeconds := cfg.App.SendIntervalSecs

	err := h.messageService.StartAutoSender(context.Background(), intervalSeconds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto sender started"})
}

func (h *MessageHandler) StopAutoSender(c *gin.Context) {
	err := h.messageService.StopAutoSender(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto sender stopped"})
}

func (h *MessageHandler) GetSentMessages(c *gin.Context) {
	messages, err := h.messageService.GetSentMessages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sent messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
