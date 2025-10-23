package repositories

import (
	"log"
	"database/sql"
	"go-project/models"
	"go-project/utils"
)


func AddDataToRag(db *sql.DB, data *models.RagData) error {
	embeddedData := utils.VectorToString(data.Embedding)

	tokens := utils.Tokenize(data.Content)
	tokensStr := utils.TokensToString(tokens)
	contentLength := len(tokens)

	query := `INSERT INTO rag_data (content, content_name, embedding, tokens, content_length)
	VALUES (?, ?, ?, ?, ?) RETURNING id
	`
	log.Printf("executing query: %s", query)
	err := db.QueryRow(query, data.Content, data.ContentName, embeddedData, tokensStr, contentLength).Scan(&data.ID)
	if err != nil {
		log.Printf("Database error: %v", err)
	}

	if err = utils.UpdateBM25Stats(db)
	err != nil {
		log.Printf("warning: failed to update BM25 stats: %v", err)
	}

	log.Printf("inserted data with %d tokens", contentLength)
	return err
}

