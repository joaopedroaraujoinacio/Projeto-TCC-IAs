package services

import (
	"fmt"
	"go-project/models"
	"go-project/utils"
)


func (s *RagService) CreateDataEmbedding(userID int64, data *models.RagData) error {
	embedding, err := utils.GenerateEmbedding(data.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}
	data.Embedding = embedding
	return s.repo.AddData(userID, data)
}

func (s *RagService) SearchSimilar(userID int64, query string, limit int) ([]models.RagData, error) {
	embedding, err := utils.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}
	return s.repo.SearchSimilar(userID, embedding, limit)
}

func (s *RagService) GetAllData(userID int64) ([]models.RagData, error) {
	return s.repo.GetAllData(userID)
}
