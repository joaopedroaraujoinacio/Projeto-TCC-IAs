package  handlers 

import (
	"net/http"
	"database/sql"
	"go-project/models"
	"go-project/services"
	"github.com/gin-gonic/gin"
)


type RagChatHandler struct {
	ragChatService services.RagChatService
	db             *sql.DB
}

func NewRagChatHandler(ragChatService services.RagChatService, db *sql.DB) *RagChatHandler {
	return &RagChatHandler{
		ragChatService: ragChatService,
		db:             db,
	}
}

func (h *RagChatHandler) RagChat(c *gin.Context) {
	var request models.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.ragChatService.RagChatService(h.db, request.Message, 1, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

