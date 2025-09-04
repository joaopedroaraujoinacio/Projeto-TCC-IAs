package routes

import (
	"database/sql"
	"golang_crud/utils"
	"golang_crud/handlers"
	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine, db *sql.DB) {

	documentHandlers := &handlers.DocumentHandlers{DB: db}
	codeHandlers := &handlers.CodeDocumentHandlers{DB: db}

	api := r.Group("/api")
	{
		api.POST("/documents", handlers.CreateDocument(db))
		api.GET("/documents/search", documentHandlers.SearchDocuments)
		api.GET("/documents", documentHandlers.GetAllDocuments)

		api.POST("/post_code", handlers.CreateCodeDocument(db))
		api.GET("/get_code", codeHandlers.SearchCodeDocuments)
		api.GET("/get_all_code", codeHandlers.GetAllCodeDocuments)

		api.GET("/hello", utils.SwaggerRoutes)
	}
}

