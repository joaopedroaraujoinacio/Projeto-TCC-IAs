package handlers

import (
	"log"
	"database/sql"
	"go-project/models"
	"github.com/gin-gonic/gin"
)


func GetAllRagData(db *sql.DB, c *gin.Context ) {
	query := "SELECT id, content FROM rag_data"

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("error: %v", err)
		c.JSON(500, gin.H{"error": "failed to fetch data on duckdb"})
		return
	}
	defer rows.Close()

	var rag_data []models.RagData
		for rows.Next() {
		var data models.RagData
	
		err := rows.Scan(&data.ID, &data.Content)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
	}

	log.Printf("found data ID %d", data.ID)
	rag_data = append(rag_data, data)

	}	
	c.JSON(200, gin.H{"rag_data": rag_data, "count": len(rag_data)})
}

