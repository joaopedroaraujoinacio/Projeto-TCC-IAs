package handlers

import (
	"io"
	"fmt"
	"log"
	"bytes"
	"net/http"
	"database/sql"
	"go-project/models"
	"go-project/services"
	"github.com/gin-gonic/gin"
)


func CreateRagData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.RagData

		body, err := c.GetRawData();
		if err != nil {
			log.Printf("failed to retreive raw data: %v", err)
			return 
		}
		fmt.Printf("raw body: %s\n", string(body))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := c.ShouldBindJSON(&data);
		err != nil {
			log.Printf("binding error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := services.CreateDataEmbedding(db, &data);
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document" })
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "embedded data created successfully",
			"embedded_data": data,
		})
	}
}

