package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
	"net/http"
)

type UtilityHandler struct {
	utilityService ports.UtilityService
}

func NewUtilityHandler(utilityService ports.UtilityService) *UtilityHandler {
	return &UtilityHandler{
		utilityService: utilityService,
	}
}

// Ping godoc
// @Summary Health check
// @Description Returns a simple pong string
// @Tags Utility
// @Success 200 {object} SuccessResponse
// @Router /ping [get]
func (h *UtilityHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// SeedSampleMessages godoc
// @Summary Seed sample messages
// @Description Seeds 10 sample messages into database for testing purposes
// @Tags Utility
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} FailResponse
// @Router /seed [post]
func (h *UtilityHandler) SeedSampleMessages(c *gin.Context) {
	err := h.utilityService.SeedSampleMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sample messages seeded successfully"})
}

// ClearDatabase godoc
// @Summary Clear database
// @Description Clears database for testing purposes
// @Tags Utility
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} FailResponse
// @Router /clear [delete]
func (h *UtilityHandler) ClearDatabase(c *gin.Context) {
	err := h.utilityService.ClearDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database cleared successfully"})
}
