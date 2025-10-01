package routes

import (
	"go-project/handlers"
	"go-project/services"
	"go-project/repositories"
	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine) {

	chatRepo := repositories.NewChatRepository("http://ollama:11434")
	chatService := services.NewChatService(chatRepo)
	chatHandler := handlers.NewChatHandler(chatService)


	api := r.Group("/api") 
	{

		api.POST("/chat", chatHandler.Chat)
	}

		r.StaticFile("/chat", "./templates/chat.html")
}

