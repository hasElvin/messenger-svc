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

// StartAutoSender godoc
// @Summary Start auto-sender
// @Description Starts the automatic message sending process
// @Tags AutoSender
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} FailResponse
// @Router /start [post]
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

// StopAutoSender godoc
// @Summary Stop auto-sender
// @Description Stops the automatic message sending process
// @Tags AutoSender
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} FailResponse
// @Router /stop [post]
func (h *MessageHandler) StopAutoSender(c *gin.Context) {
	err := h.messageService.StopAutoSender(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto sender stopped"})
}

// GetSentMessages godoc
// @Summary List all sent messages
// @Description Returns a list of messages that were sent by the auto-sender
// @Tags Messages
// @Success 200 {array} domain.Message
// @Router /sent [get]
func (h *MessageHandler) GetSentMessages(c *gin.Context) {
	messages, err := h.messageService.GetSentMessages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sent messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
