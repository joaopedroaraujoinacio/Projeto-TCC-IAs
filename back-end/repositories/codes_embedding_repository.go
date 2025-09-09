package repositories

import (
	"log"
	"database/sql"
	"golang_crud/utils"
	"golang_crud/models"
)


func CreateCodeDocument(db *sql.DB, code *models.Codes) error {
	embeddingStr := utils.VectorToString(code.Embedding)	

	query := `INSERT INTO codes (content, programming_language, file_name, embedding) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	log.Printf("Executing query: %s", query)
	log.Printf("Parameters: content=%s, programming_language=%s, file_name=%s, embedding=%d", code.Content, code.ProgramingLanguage, code.FileName, len(code.Embedding))

	err := db.QueryRow(query, code.Content, code.ProgramingLanguage, code.FileName, embeddingStr).Scan(&code.ID, &code.CreatedAt)
	if err != nil {
		log.Printf("Database error: %v", err)
	}
	return err
}


func SearchSimilarCodeDocuments(db *sql.DB, queryEmbedding []float32, limit int) ([]models.Codes, error) {
 embeddingStr := utils.VectorToString(queryEmbedding)

	query := `
		SELECT id, content, programming_language, file_name, embedding, created_at,
		       1 - (embedding <=> $1) as similarity
		FROM codes 
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

	var codes []models.Codes
	for rows.Next() {
		var code models.Codes
		var similarity float64
		var embeddingStr string
		var programmingLanguage, fileName sql.NullString

		err := rows.Scan(&code.ID, &code.Content, &programmingLanguage, &fileName, &embeddingStr, &code.CreatedAt, &similarity)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		if programmingLanguage.Valid {
			code.ProgramingLanguage = &programmingLanguage.String
		} 

		if fileName.Valid {
			code.FileName = &fileName.String
		}

		code.Embedding = utils.ParseVectorString(embeddingStr) 

		log.Printf("Found code ID %d with similarity: %.4f", code.ID, similarity)
		codes = append(codes, code)
		}
		return codes, nil
	}

