package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"go-project/models"
	"github.com/gin-gonic/gin"
)

func (h *ChatHandler) StreamGemini(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	messageChan, errorChan := h.chatService.StreamGemini(&request)

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
			if strings.HasPrefix(msg, "__token_stats__:") {
				fmt.Fprintf(w, "event: token_stats\ndata: %s\n\n", strings.TrimPrefix(msg, "__token_stats__:"))
				return true
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
