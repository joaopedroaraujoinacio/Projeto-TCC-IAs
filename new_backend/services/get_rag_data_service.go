package services

import (
	"fmt"
	"log"
	"database/sql"
	"go-project/utils"
	"go-project/models"
	"go-project/repositories"
)


func SearchSimilarData(db *sql.DB, query string, limit int, bm25Weight float64, vectorWeight float64) ([]models.RagData, error) {
	// Validate weights
	if bm25Weight < 0 || vectorWeight < 0 {
		return nil, fmt.Errorf("weights must be non-negative")
	}
	
	// If both weights are 0, set default balanced weights
	if bm25Weight == 0 && vectorWeight == 0 {
		bm25Weight = 0.5
		vectorWeight = 0.5
	}
	
	// Generate embedding for the query
	queryEmbedding, err := utils.GenerateEmbedding(query)
	if err != nil {
		log.Printf("failed to generate query embedding: %v", err)
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}
	
	log.Printf("Generated embedding for query: '%s' (length: %d)", query, len(queryEmbedding))
	if len(queryEmbedding) > 0 {
		log.Printf("first 5 embedding values: %v", queryEmbedding[:5])
	}
	
	// Perform hybrid search
	results, err := repositories.RagDataSimilaritySearch(db, query, queryEmbedding, limit, bm25Weight, vectorWeight)
	if err != nil {
		return nil, err
	}
	
	log.Printf("repository returned %d results", len(results))
	return results, nil
}
// func SearchSimilarData(db *sql.DB, query string, limit int) ([]models.RagData, error) {
// 	queryEmbedding, err := utils.GenerateEmbedding(query)
// 	if err != nil {
// 		log.Printf("failed to generate query embedding: %v", err)
// 		return nil, fmt.Errorf("failed to generate query embedding :%w", err)
// 	}
// 	log.Printf("Generated embedding for query: '%s' (lenght: %d)", query, len(queryEmbedding))
//
// 	if len(queryEmbedding) > 0 {
// 		log.Printf("first 5 embedding values: %v", queryEmbedding[:5])
// 	}
//
// 	results, err := repositories.RagDataSimilaritySearch(db, queryEmbedding, limit)
//
// 	log.Printf("repository returned %d results", len(results) )
//
// 	return results, err
// }
//
