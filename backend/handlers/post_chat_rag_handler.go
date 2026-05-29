package handlers

import (
	"fmt"
	"go-project/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *ChatHandler) RagChat(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	messageChan, errorChan := h.chatService.RagChat(userID, &request)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				select {
				case err, ok := <-errorChan:
					if ok && err != nil {
						fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
						return false
					}
				default:
				}
				fmt.Fprintf(w, "event: done\ndata: \n\n")
				return false
			}
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
			return true
		case err, ok := <-errorChan:
			if ok && err != nil {
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			}
			return false
		case <-c.Request.Context().Done():
			return false
		}
	})
}

func getUserIDFromContext(c *gin.Context) (int64, bool) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, false
	}
	id, ok := val.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return 0, false
	}
	return id, true
}
