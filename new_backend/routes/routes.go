package routes

import (
	"database/sql"
	"go-project/handlers"
	"go-project/services"
	"go-project/repositories"
	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine, db *sql.DB) {

	chatRepo := repositories.NewChatRepository("http://ollama:11434")
	chatService := services.NewChatService(chatRepo)
	chatHandler := handlers.NewChatHandler(chatService)

	chatStreamService := services.NewChatStreamService(chatRepo)
	chatStreamHandler := handlers.NewChatStreamHandler(chatStreamService)

	ragChatService := services.NewRagChatService(chatRepo)
	ragChatHandler := handlers.NewRagChatHandler(ragChatService, db)

	api := r.Group("/api") 
	{
		api.POST("/chat", chatHandler.Chat)
		api.POST("/rag/chat", ragChatHandler.RagChat)
		api.POST("/rag", handlers.CreateRagData(db))
		api.GET("/rag/search", handlers.SearchSimilarRagData(db))
		// http://localhost:8080/api/rag/search?q=

		api.POST("/chat/stream", chatStreamHandler.StreamChat)
	}

		r.StaticFile("/chat", "./templates/chat.html")
}

