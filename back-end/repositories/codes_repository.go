package repositories

import (
	"database/sql"
	"golang_crud/models"
	"golang_crud/utils"
	"log"
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
