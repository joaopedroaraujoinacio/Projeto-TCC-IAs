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

	ragChatService := services.NewRagChatService(chatRepo)
	ragChatHandler := handlers.NewRagChatHandler(ragChatService, db)

	api := r.Group("/api") 
	{
		api.POST("/chat", chatHandler.StreamChat)
		api.POST("/chat/rag", ragChatHandler.RagChat)
		api.POST("/rag", handlers.CreateRagData(db))
		api.GET("/rag/search", handlers.SearchSimilarRagData(db))
	}

		r.StaticFile("/chat", "./templates/chat.html")
}

