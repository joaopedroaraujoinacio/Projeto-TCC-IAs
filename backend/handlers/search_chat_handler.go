package handlers

import (
	"net/http"
	"golang_crud/models"
	"golang_crud/services"
	"github.com/gin-gonic/gin"
)

type SmartSearchHandler struct {
	smartSearchService services.SmartSearchService
}

func NewSmartSearchHandler(smartSearchService services.SmartSearchService) *SmartSearchHandler {
	return &SmartSearchHandler{
		smartSearchService: smartSearchService,
	}
}

func (h *SmartSearchHandler) SmartSearch(c *gin.Context) {
	var request models.SmartSearchRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	response, err := h.smartSearchService.ProcessSmartSearch(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process smart search",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *SmartSearchHandler) SmartSearchHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "Smart Search service is running",
		"provider": "DuckDuckGo via Tor",
		"ai_model": "llama3.2:3b",
	})
}

