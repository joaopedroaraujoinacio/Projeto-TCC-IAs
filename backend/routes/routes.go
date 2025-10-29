package routes

import (
	"database/sql"
	"go-project/handlers"
	"go-project/repositories"
	"go-project/services"
	"go-project/utils"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine, db *sql.DB) {

	chatRepo := repositories.NewChatRepository("http://ollama:11434")
	searchRepo := utils.NewWebSearchRepository()

	chatService := services.NewChatService(chatRepo, searchRepo)
	
	chatHandler := handlers.NewChatHandler(chatService)
	ragChatHandler := handlers.NewRagChatHandler(chatService, db)
	webSearchHandler := handlers.NewWebSearchHandler(chatService)

	getAllRagData := func(c *gin.Context) {
		handlers.GetAllRagData(db, c)
	}

	api := r.Group("/api") 
	{
		api.POST("/chat", chatHandler.StreamChat)
		api.POST("/chat/rag", ragChatHandler.RagChat)
		api.POST("/rag", handlers.CreateRagData(db))
		api.POST("/chat/web-search", webSearchHandler.WebSearchChat)
		api.GET("/rag/search", handlers.SearchSimilarRagData(db))
		api.GET("/rag/get_all_data", getAllRagData)
		
	}

		r.StaticFile("/chat", "./templates/chat.html")
}

