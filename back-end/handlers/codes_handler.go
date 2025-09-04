package handlers

import (
	"io"
	"fmt"
	"log"
	"bytes"
	"strconv"
	"net/http"
	"database/sql"
	"golang_crud/models"
	"golang_crud/services"
	"github.com/gin-gonic/gin"
)


func CreateCodeDocument(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		var code models.Codes

		body, _ := c.GetRawData()
		fmt.Printf("Raw body: %s\n", string(body))

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := c.ShouldBindJSON(&code); err != nil {
			fmt.Printf("Binding error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		fmt.Printf("Parsed document: %+v\n", code)

		if err := services.CreateCodeDocument(db, &code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Faled to create document"})
			return 
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"message": "Code document created successfully",
			"code": code,
		})	
	}
}


type CodeDocumentHandlers struct {
	DB *sql.DB
}

func (h *CodeDocumentHandlers) SearchCodeDocuments(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	log.Printf("Searching for: %s (limit: %d)", query, limit)

	codes, err := services.SearchCodeDocuments(h.DB, query, limit)
	if err != nil {
		log.Printf("Search error: %v", err)
		c.JSON(500, gin.H{"error": "Failed to search documents"})
		return
	}

	log.Printf("Found %d similar documents", len(codes))
	c.JSON(200, gin.H{
		"query": query,
		"results": codes,
		"count": len(codes),
	})
}


func (h *CodeDocumentHandlers) GetAllCodeDocuments(c *gin.Context) {
	query := `SELECT id, content, programming_language, file_name, created_at FROM codes ORDER BY id DESC LIMIT 10`
	
	rows, err := h.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()
	
	var codes []models.Codes
	for rows.Next() {
		var code models.Codes
		rows.Scan(&code.ID, &code.Content, &code.ProgramingLanguage, &code.FileName, &code.CreatedAt)
		codes = append(codes, code)
	}
	
	c.JSON(200, gin.H{"codes": codes, "count": len(codes)})
}

