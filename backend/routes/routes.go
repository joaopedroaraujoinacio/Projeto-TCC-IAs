package routes

import (
	"database/sql"
	"golang_crud/utils"
	"golang_crud/handlers"
	"golang_crud/services"
	"golang_crud/repositories"
	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine, db *sql.DB) {

	documentHandlers := &handlers.DocumentHandlers{DB: db}
	codeHandlers := &handlers.CodeDocumentHandlers{DB: db}

	chatRepo := repositories.NewChatRepository("http://ollama:11434")
	chatService := services.NewChatService(chatRepo)
	chatHandler := handlers.NewChatHandler(chatService)

	smartSearchRepo := repositories.NewSmartSearchRepository()
	smartSearchService := services.NewSmartSearchService(smartSearchRepo, chatRepo)
	smartSearchHandler := handlers.NewSmartSearchHandler(smartSearchService)

	api := r.Group("/api")
	{
		api.POST("/documents", handlers.CreateDocument(db))
		api.GET("/documents/search", documentHandlers.SearchDocuments)
		api.GET("/documents", documentHandlers.GetAllDocuments)

		api.POST("/post_code", handlers.CreateCodeDocument(db))
		api.GET("/get_code", codeHandlers.SearchCodeDocuments)
		api.GET("/get_all_code", codeHandlers.GetAllCodeDocuments)

		api.POST("/chat", chatHandler.Chat)
		api.GET("/hello", utils.SwaggerRoutes)

		api.POST("/smart-search", smartSearchHandler.SmartSearch)
		api.GET("/smart-search/health", smartSearchHandler.SmartSearchHealth)

	}
}

