package repositories

import (
	"log"
	"database/sql"
	"golang_crud/models"
	"golang_crud/utils"
)


func CreateDocument(db *sql.DB, doc *models.Document) error {
	embeddingStr := utils.VectorToString(doc.Embedding)
	
	query := `INSERT INTO documents (content, media_type, file_name, embedding) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	
	log.Printf("Executing query: %s", query)
	log.Printf("Parameters: content=%s, media_type=%s, file_name=%s, embedding_length=%d", 
		doc.Content, doc.MediaType, doc.FileName, len(doc.Embedding))
	
	err := db.QueryRow(query, doc.Content, doc.MediaType, doc.FileName, embeddingStr).Scan(&doc.ID, &doc.CreatedAt)
	if err != nil {
		log.Printf("Database error: %v", err)
	}
	return err
}


func SearchSimilarDocuments(db *sql.DB, queryEmbedding []float32, limit int) ([]models.Document, error) {
	embeddingStr := utils.VectorToString(queryEmbedding)

	query := `
		SELECT id, content, media_type, file_name, embedding, created_at,
		       1 - (embedding <=> $1) as similarity
		FROM documents 
		ORDER BY embedding <=> $1 
		LIMIT $2
	`

	log.Printf("Executing search query with embedding length: %d, limit: %d", len(queryEmbedding), limit)

	rows, err := db.Query(query, embeddingStr, limit)
	if err != nil {
		log.Printf("Search query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		var similarity float64
		var embeddingStr string
		var mediaType, fileName sql.NullString

		err := rows.Scan(&doc.ID, &doc.Content, &mediaType, &fileName, &embeddingStr, &doc.CreatedAt, &similarity)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		if mediaType.Valid {
			doc.MediaType = &mediaType.String
		}
		if fileName.Valid {
			doc.FileName = &fileName.String
		}
		
		doc.Embedding = utils.ParseVectorString(embeddingStr)

		log.Printf("Found document ID %d with similarity: %.4f", doc.ID, similarity)
		documents = append(documents, doc)
	}

	return documents, nil
}

