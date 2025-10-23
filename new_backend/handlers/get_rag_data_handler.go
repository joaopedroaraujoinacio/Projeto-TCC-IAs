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
		
		// Get limit parameter
		limitStr := c.DefaultQuery("limit", "5")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 5
		}
		
		// Get weight parameters (optional, defaults to 0.5 each for balanced search)
		bm25WeightStr := c.DefaultQuery("bm25_weight", "0.5")
		bm25Weight, err := strconv.ParseFloat(bm25WeightStr, 64)
		if err != nil {
			bm25Weight = 0.5
		}
		
		vectorWeightStr := c.DefaultQuery("vector_weight", "0.5")
		vectorWeight, err := strconv.ParseFloat(vectorWeightStr, 64)
		if err != nil {
			vectorWeight = 0.5
		}
		
		log.Printf("Hybrid search for: %s (limit: %d, bm25: %.2f, vector: %.2f)", 
			query, limit, bm25Weight, vectorWeight)
		
		data, err := services.SearchSimilarData(db, query, limit, bm25Weight, vectorWeight)
		if err != nil {
			log.Printf("search error: %v", err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("failed to search documents: %v", err)})
			return
		}
		
		log.Printf("found %d similar documents", len(data))
		c.JSON(200, gin.H{
			"query":         query,
			"results":       data,
			"bm25_weight":   bm25Weight,
			"vector_weight": vectorWeight,
		})
	}
}

// func SearchSimilarRagData(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
//
// 	query := c.Query("q")
// 	if query == "" {
// 		c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
// 		return
// 	}
//
// 	limitStr := c.DefaultQuery("limit", "1")
// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil {
// 		limit = 1
// 	}
//
// 	log.Printf("searching for: %s (limit: %d)", query, limit)
//
// 	data, err := services.SearchSimilarData(db, query, limit)
// 	if err != nil {
// 		log.Printf("search error: %v", err)
// 		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to search documents: %v", err)})
// 		return
// 	}
//
// 	log.Printf("found %d similar documents", len(data))
// 	c.JSON(200, gin.H{
// 		"query": query,
// 		"results": data,
// 		})
// 	}
// }
