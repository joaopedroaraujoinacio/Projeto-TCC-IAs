package services

import (
	"fmt"
	"log"
	"database/sql"
	"go-project/utils"
	"go-project/models"
	"go-project/repositories"
)


func SearchSimilarData(db *sql.DB, query string, limit int) ([]models.RagData, error) {
	queryEmbedding, err := utils.GenerateEmbedding(query)
	if err != nil {
		log.Printf("failed to generate query embedding: %v", err)
		return nil, fmt.Errorf("failed to generate query embedding :%w", err)
	}
	log.Printf("Generated embedding for query: '%s' (lenght: %d)", query, len(queryEmbedding))

	if len(queryEmbedding) > 0 {
		log.Printf("first 5 embedding values: %v", queryEmbedding[:5])
	}

	results, err := repositories.RagDataSimilaritySearch(db, queryEmbedding, limit)

	log.Printf("repository returned %d results", len(results) )

	return results, err
}

