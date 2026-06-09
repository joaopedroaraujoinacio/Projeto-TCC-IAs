package routes

import (
	"go-project/handlers"
	"go-project/middleware"

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
		api.POST("/auth/login", deps.UserHandler.Login)

		api.POST("/chat", deps.ChatHandler.StreamChat)
		api.POST("/chat/web-search", deps.ChatHandler.WebSearchChat)
		api.POST("/chat/openai", deps.ChatHandler.StreamOpenAI)
		api.POST("/chat/gemini", deps.ChatHandler.StreamGemini)

		authenticated := api.Group("/rag")
		authenticated.Use(middleware.Authenticate())
		{
		authenticated.POST("/rag_chat", deps.ChatHandler.RagChat)
		authenticated.POST("/add_rag_data", deps.RagHandler.CreateRagData)
		authenticated.GET("/search_rag_data", deps.RagHandler.SearchSimilarRagData)
		authenticated.GET("/get_all_rag_data", deps.RagHandler.GetAllRagData)
		}
	}
		server.StaticFile("/chat", "./templates/index.html")
		server.StaticFile("/login", "./templates/login.html")
		server.StaticFile("/signup", "./templates/signup.html")
}

