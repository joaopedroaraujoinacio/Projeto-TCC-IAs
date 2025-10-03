package handlers

import (
	"fmt"
	"log"
	"strconv"
	"database/sql"
	"go-project/services"
	"github.com/gin-gonic/gin"
)


func SearchSimilarRagData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 1
	}

	log.Printf("searching for: %s (limit: %d)", query, limit)

	data, err := services.SearchSimilarData(db, query, limit)
	if err != nil {
		log.Printf("search error: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to search documents: %v", err)})
		return
	}

	log.Printf("found %d similar documents", len(data))
	c.JSON(200, gin.H{
		"query": query,
		"results": data,
		})
	}
}
