package handlers

import (
	"go-project/models"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *ChatHandler) StreamOpenAI(c *gin.Context) {
    var request models.ChatRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("Access-Control-Allow-Origin", "*")

    messageChan, errorChan := h.chatService.StreamOpenAI(&request)

    c.Stream(func(w io.Writer) bool {
        select {
        case msg, ok := <-messageChan:
            if !ok {
                c.SSEvent("done", "")
                return false
            }
            if strings.HasPrefix(msg, "__token_stats__:") {
                c.SSEvent("token_stats", strings.TrimPrefix(msg, "__token_stats__:"))
                return true
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
