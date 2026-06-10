package handlers

import (
	"go-project/models"
	"go-project/utils"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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
	
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".txt": true,
		",md": true,
		".csv": true,
		".pdf": true,
	}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsuported file type."})
		return
	}

	var content string

	if ext == ".pdf" {
		content, err = utils.ExtractPDFText(file, header.Size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

	} else {
		raw, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}
		content = string(raw)
	}

	content = strings.TrimSpace(content)
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file appears to be empty"})
		return
	}

	data := models.RagData{
		Content: content,
		ContentName: header.Filename,
	}

	if err := h.svc.CreateDataEmbedding(userID, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to index document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": 			"document indexed successfully",
		"content_name": data.ContentName,
		"id": 					data.ID,
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

