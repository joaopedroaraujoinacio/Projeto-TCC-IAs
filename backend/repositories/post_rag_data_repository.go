package repositories

import (
	"fmt"
	"log"
	"go-project/models"
	"go-project/utils"
)

func (r *ragRepositoryImpl) AddData(userID int64, data *models.RagData) error {
	embeddedData := utils.VectorToString(data.Embedding)
	tokens := utils.Tokenize(data.Content)
	tokensStr := utils.TokensToString(tokens)
	contentLength := len(tokens)

	query := `
		INSERT INTO rag_data (user_id, content, content_name, embedding, tokens, content_length)
		VALUES (?, ?, ?, ?, ?, ?) RETURNING id
	`
	err := r.db.QueryRow(query, userID, data.Content, data.ContentName, embeddedData, tokensStr, contentLength).Scan(&data.ID)
	if err != nil {
		return fmt.Errorf("failed to insert rag data: %w", err)
	}
	if err := utils.UpdateBM25Stats(r.db, userID); err != nil {
		log.Printf("warning: failed to update BM25 stats: %v", err)
	}
	log.Printf("inserted rag data with %d tokens for user %d", contentLength, userID)
	return nil
}

func (r *ragRepositoryImpl) SearchSimilar(userID int64, embedding []float32, limit int) ([]models.RagData, error) {
	embeddedData := utils.VectorToString(embedding)
	query := `
		SELECT id, content, content_name
		FROM rag_data
		WHERE user_id = ?
		ORDER BY array_distance(embedding, ?::FLOAT[768])
		LIMIT ?
	`
	rows, err := r.db.Query(query, userID, embeddedData, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search rag data: %w", err)
	}
	defer rows.Close()

	var results []models.RagData
	for rows.Next() {
		var d models.RagData
		if err := rows.Scan(&d.ID, &d.Content, &d.ContentName); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	return results, nil
}

func (r *ragRepositoryImpl) GetAllData(userID int64) ([]models.RagData, error) {
	query := `SELECT id, content, content_name FROM rag_data WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rag data: %w", err)
	}
	defer rows.Close()

	var results []models.RagData
	for rows.Next() {
		var d models.RagData
		if err := rows.Scan(&d.ID, &d.Content, &d.ContentName); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	return results, nil
}
