package routes

import (
	"net/http"
	"database/sql"
	"golang_crud/handlers"
	"github.com/gin-gonic/gin"
)


func Setup(r *gin.Engine, db *sql.DB) {

	documentHandlers := &handlers.DocumentHandlers{DB: db}

	api := r.Group("/api")
	{
		api.POST("/documents", handlers.CreateDocument(db))
		api.GET("/documents/search", documentHandlers.SearchDocuments)
		api.GET("/documents", documentHandlers.GetAllDocuments)

		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello world",
			})
		})
	}
}

