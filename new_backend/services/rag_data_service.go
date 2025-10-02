package services

import (
	"fmt"
	"database/sql"
	"go-project/utils"
	"go-project/models"
	"go-project/repositories"
)


func CreateDataEmbedding(db *sql.DB, data *models.RagData) error {
	embedding, err := utils.GenerateEmbedding(data.Content) 
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	data.Embedding = embedding
	return repositories.AddDataToRag(db, data)

}

