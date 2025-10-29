package handlers

import (
	"go-project/models"
	"go-project/services"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


type WebSearchHandler struct {
	webSearchService services.ChatService
}

func NewWebSearchHandler(webSearchService services.ChatService) *WebSearchHandler {
	return &WebSearchHandler{
		webSearchService: webSearchService,
	}
}

func (h *WebSearchHandler) WebSearchChat(c *gin.Context) {
	var request models.WebSearchRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	response, messageChan, errorChan := h.webSearchService.WebSearchChat(&request)

	if response != nil {
		log.Printf("Using %d sources:", response.Count)
		for i, source := range response.Sources {
			log.Printf("  [%d] %s - %s", i+1, source.Title, source.URL)
		}
		
		c.SSEvent("sources", response)
		c.Writer.Flush()
	}

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

