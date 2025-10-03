package repositories

import (
	"log"
	"database/sql"
	"go-project/utils"
	"go-project/models"
)


func RagDataSimilaritySearch(db *sql.DB, queryEmbedding []float32, limit int) ([]models.RagData, error) {
	embeddedData := utils.VectorToString(queryEmbedding)

	log.Printf("embedding string first 100 chars: %s", embeddedData[:100])

	query := `
    SELECT id, content,
           array_cosine_similarity(embedding, CAST(? AS FLOAT[768])) as similarity
    FROM rag_data 
    ORDER BY similarity DESC
		LIMIT ?
`
	log.Printf("Executing search query with embedding length: %d, limit: %d", len(queryEmbedding), limit)

	rows, err := db.Query(query, embeddedData, limit)
	if err != nil {
		log.Printf("search query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	var rag_data []models.RagData
		for rows.Next() {
		var data models.RagData
		var similarity float64
	
		err := rows.Scan(&data.ID, &data.Content, &similarity)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
	}

	log.Printf("found data ID %d with similarity: %.4f", data.ID, similarity)
	rag_data = append(rag_data, data)

	}	

	return rag_data, nil
} 

