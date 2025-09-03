package services

import (
	"fmt"
	"log"
	"database/sql"
	"golang_crud/utils"
	"golang_crud/models"
	"golang_crud/repositories"
)


func CreateCodeDocument(db *sql.DB, code *models.Codes) error {
	embedding, err := utils.GenerateCodeEmbedding(code.Content) 
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	code.Embedding = embedding
	return repositories.CreateCodeDocument(db, code)
}


func SearchCodeDocuments(db *sql.DB, query string, limit int) ([]models.Codes, error) {
	queryEmbedding, err := utils.GenerateEmbedding(query)
	if err != nil {
		log.Printf("Failed to generate embedding for query: %v", err)
		return nil, fmt.Errorf("failed to generate embedding for query: %w", err)
	}

	log.Printf("Generated embedding for query: '%s' (lenght: %d)", query, len(queryEmbedding))
	
	return repositories.SearchSimilarCodeDocuments(db, queryEmbedding, limit)
}
