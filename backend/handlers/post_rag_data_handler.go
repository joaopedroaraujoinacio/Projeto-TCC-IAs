package handlers

import (
	"net/http"
	"strconv"
	"go-project/models"
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (int64, bool) {
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

func (h *RagHandler) CreateRagData(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var data models.RagData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateDataEmbedding(userID, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":       "embedded data created successfully",
		"embedded_data": data,
	})
}

func (h *RagHandler) SearchSimilarRagData(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	results, err := h.svc.SearchSimilar(userID, query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *RagHandler) GetAllRagData(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	results, err := h.svc.GetAllData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}
