package handlers

import (
	"net/http"
	"golang_crud/models"
	"golang_crud/services"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService services.ChatService
}

func NewChatHandler(chatService services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var request models.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	response, err := h.chatService.ProcessChat(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process chat",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *ChatHandler) ChatHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Chat service is running",
		"model":  "llama3.2:3b",
	})
}

