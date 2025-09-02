package services

import (
	"fmt"
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

