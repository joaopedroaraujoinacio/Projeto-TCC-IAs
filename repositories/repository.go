package repositories

import (
	"database/sql"
	"fmt"
	"golang_crud/models"
	"log"
	"strconv"
	"strings"
)


func CreateDocument(db *sql.DB, doc *models.Document) error {
	embeddingStr := vectorToString(doc.Embedding)
	
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

// Helper function to convert []float32 to pgvector string format
func vectorToString(embedding []float32) string {
	strValues := make([]string, len(embedding))
	for i, v := range embedding {
		strValues[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(strValues, ",") + "]"
}



// func CreateDocument(db *sql.DB, doc *models.Document) error {
// 	query := `INSERT INTO documents (content, media_type, file_name, embedding) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
// 	err := db.QueryRow(query, doc.Content, doc.MediaType, doc.FileName, pq.Array(doc.Embedding)).Scan(&doc.ID, &doc.CreatedAt)
// 	return err
// }


// func CreateDocument(db *sql.DB, doc *models.Document) error {
//   query := `INSERT INTO documents (content, embedding) VALUES ($1, $2) RETURNING id, created_at`
// 	err := db.QueryRow(query, doc.Content, pq.Array(doc.Embedding)).Scan(&doc.ID, &doc.CreatedAt)
// 	return err
// }

func SearchSimilarDocuments(db *sql.DB, queryEmbedding []float32, limit int) ([]models.Document, error) {
	embeddingStr := vectorToString(queryEmbedding)

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
		
		doc.Embedding = parseVectorString(embeddingStr)

		log.Printf("Found document ID %d with similarity: %.4f", doc.ID, similarity)
		documents = append(documents, doc)
	}

	return documents, nil
}

func parseVectorString(vectorSTR string) []float32 {
	vectorSTR = strings.Trim(vectorSTR, "[]")
	if vectorSTR == "" {
		return []float32{}
	}

	parts := strings.Split(vectorSTR, ",")
	embedding := make([]float32, len(parts))
	
	for i, part := range  parts {
		if val, err := strconv.ParseFloat(strings.TrimSpace(part), 32); err == nil {
			embedding[i] = float32(val)
		}
	}

	return embedding
}
