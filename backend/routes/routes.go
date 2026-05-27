package routes

import (
	"go-project/handlers"

	"github.com/gin-gonic/gin"
)


type Dependencies struct {
	ChatHandler *handlers.ChatHandler
	UserHandler *handlers.UserHandler
	RagHandler 	*handlers.RagHandler
}

func SetupRoutes(server *gin.Engine, deps *Dependencies) {
	api := server.Group("/api") 
	{

		api.POST("/auth/signup", deps.UserHandler.Signup)
		api.POST("auth/login", deps.UserHandler.Login)

		api.POST("/chat", deps.ChatHandler.StreamChat)
		api.POST("/chat/rag", deps.ChatHandler.RagChat)
		api.POST("/chat/web-search", deps.ChatHandler.WebSearchChat)
		api.POST("/chat/openai", deps.ChatHandler.StreamOpenAI)
		api.POST("/chat/gemini", deps.ChatHandler.StreamGemini)

		api.POST("/rag", deps.RagHandler.CreateRagData)
		api.GET("/rag/search", deps.RagHandler.SearchSimilarRagData)
		api.GET("/rag/get_all_data", deps.RagHandler.GetAllRagData)
	}

		server.StaticFile("/chat", "./templates/index.html")
}

