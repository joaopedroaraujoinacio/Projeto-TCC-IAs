package handlers

import (
	"io"
	"fmt"
	"bytes"
	"net/http"
	"database/sql"
	"golang_crud/models"
	"golang_crud/services"
	"github.com/gin-gonic/gin"
)


func CreateCodeDocument(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		var code models.Codes

		body, _ := c.GetRawData()
		fmt.Printf("Raw body: %s\n", string(body))

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := c.ShouldBindJSON(&code); err != nil {
			fmt.Printf("Binding error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		fmt.Printf("Parsed document: %+v\n", code)

		if err := services.CreateCodeDocument(db, &code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Faled to create document"})
			return 
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"message": "Code document created successfully",
			"code": code,
		})	
	}
}

