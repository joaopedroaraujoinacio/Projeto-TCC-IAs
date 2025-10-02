package repositories

import (
	"log"
	"database/sql"
	"go-project/models"
	"go-project/utils"
)


func AddDataToRag(db *sql.DB, data *models.RagData) error {
	embeddedData := utils.VectorToString(data.Embedding)

	query := `INSERT INTO rag_data (content, content_name, embedding)
	VALUES (?, ?, ?) RETURNING id
	`
	log.Printf("executing query: %s", query)

	err := db.QueryRow(query, data.Content, data.ContentName, embeddedData).Scan(&data.ID)
	if err != nil {
		log.Printf("Database error: %v", err)
	}
	return err
}

