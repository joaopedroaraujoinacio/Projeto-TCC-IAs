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



func (h *DocumentHandlers) SearchCodeDocuments(c *gin.Context) {
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

