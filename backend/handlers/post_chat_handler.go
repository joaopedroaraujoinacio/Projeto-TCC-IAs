package handlers

import (
	"io"
	"net/http"
	"go-project/models"
	"github.com/gin-gonic/gin"
)

// StreamChat godoc
// @Summary      Send a message to the chat
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        request body models.ChatRequest true "Chat message"
// @Success      200 {object} map[string]string "Chat response"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/chat [post]
func (h *ChatHandler) StreamChat(c *gin.Context) {
	var request models.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	messageChan, errorChan := h.chatService.StreamChat(&request)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-messageChan:
			if !ok {
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
			return false
		}
	})
}

