package services

import (
	"database/sql"
	"fmt"
	"golang_crud/models"
	"golang_crud/repositories"
	"log"
)


func CreateDocument(db *sql.DB, doc *models.Document) error {
	embedding, err := GenerateEmbedding(doc.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	doc.Embedding = embedding
	return repositories.CreateDocument(db, doc)
}

func SearchDocuments(db *sql.DB, query string, limit int) ([]models.Document, error) {
	queryEmbedding, err := GenerateEmbedding(query)
	if err != nil {
		log.Printf("Failed to generate embedding for query: %v", err)
		return nil, fmt.Errorf("failed to generate embedding for query: %w", err)
	}

	log.Printf("Generated embedding for query: '%s' (lenght: %d)", query, len(queryEmbedding))
	
	return repositories.SearchSimilarDocuments(db, queryEmbedding, limit)
}
