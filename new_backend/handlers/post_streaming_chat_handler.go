package handlers

import (
	"io"
	"net/http"
	"go-project/models"
	"go-project/services"
	"github.com/gin-gonic/gin"
)

type ChatStreamHandler struct {
	chatStreamService services.ChatStreamService
}

func NewChatStreamHandler(chatStreamService services.ChatStreamService) *ChatStreamHandler {
	return &ChatStreamHandler{
		chatStreamService: chatStreamService,
	}
}

// StreamChat handles POST /api/chat/stream
func (h *ChatStreamHandler) StreamChat(c *gin.Context) {
	var request models.ChatRequest

	// Parse request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Get streaming channels from service
	messageChan, errorChan := h.chatStreamService.StreamChat(&request)

	// Stream to client
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				// Stream complete
				c.SSEvent("done", "")
				return false
			}
			c.SSEvent("message", msg)
			return true

		case err, ok := <-errorChan:
			if ok && err != nil {
				c.SSEvent("error", err.Error())
			}
			return false

		case <-c.Request.Context().Done():
			// Client disconnected
			return false
		}
	})
}
