package repositories

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"golang_crud/models"
)


type ImageEmbeddingRepository struct {
	db *sql.DB
}

func NewImageEmbeddingRepository(db *sql.DB) *ImageEmbeddingRepository{
	return &ImageEmbeddingRepository{db: db}
}

func (r *ImageEmbeddingRepository) Create(embedding *models.ImageEmbeding) error{
	metadataJSON, err := json.Marshal(embedding.Metadata)
	if err != nil {
		return err
	}

	query := `
			INSERT INTO image_embeddings (filename, file_path, description, embedding, metadata)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at
	`

	return r.db.QueryRow(
			query,
			embedding.Filename,
			embedding.FilePath,
			embedding.Description,
			pq.Array(embedding.Embedding),
			metadataJSON,
		).Scan(&embedding.ID, &embedding.CreatedAt)
}

func (r *ImageEmbeddingRepository) FindSimilar(embedding []float32, limit int) ([]models.ImageEmbeding, error) {
    query := `
        SELECT id, filename, file_path, description, embedding, metadata, created_at,
               1 - (embedding <=> $1) as similarity
        FROM image_embeddings
        ORDER BY embedding <=> $1
        LIMIT $2
    `
		rows, err := r.db.Query(query, pq.Array(embedding), limit)
		if err != nil {
			return nil, err
	}
	defer rows.Close()

	var results []models.ImageEmbeding
	for rows.Next() {
		var img models.ImageEmbeding
		var metadataJSON []byte
		var similarity float64
		
		err := rows.Scan(
				&img.ID, &img.Filename, &img.FilePath, &img.Description,
				pq.Array(&img.Embedding), &metadataJSON, &img.CreatedAt, &similarity,
		)
		if err != nil {
				return nil, err
		}
		
		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &img.Metadata)
	}
	 results = append(results, img)
	}

	return results, nil

}

