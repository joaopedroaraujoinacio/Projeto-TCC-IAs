package routes

import (
	"go-project/handlers"

	"github.com/gin-gonic/gin"
)


type Dependencies struct {
	ChatHandler *handlers.ChatHandler
}

func SetupRoutes(server *gin.Engine, deps *Dependencies) {
	// getAllRagData := func(c *gin.Context) {
	// 	handlers.GetAllRagData(db, c)
	// }

	api := server.Group("/api") 
	{
		api.POST("/chat", deps.ChatHandler.StreamChat)
		api.POST("/chat/rag", deps.ChatHandler.RagChat)
		api.POST("/chat/web-search", deps.ChatHandler.WebSearchChat)
		api.POST("/chat/openai", deps.ChatHandler.StreamOpenAI)
		// api.POST("/rag", handlers.CreateRagData(db))
		// api.GET("/rag/search", handlers.SearchSimilarRagData(db))
		// api.GET("/rag/get_all_data", getAllRagData)
	}

		server.StaticFile("/chat", "./templates/index.html")
}

