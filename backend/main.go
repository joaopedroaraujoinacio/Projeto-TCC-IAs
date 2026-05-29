package main

import (
	"go-project/config"
	"go-project/handlers"
	"go-project/repositories"
	"go-project/routes"
	"go-project/server"
	"go-project/services"
	"go-project/utils"
	"log"
	"os"

	_ "go-project/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Farol IA
// @version         1.0
// @description     Self hosted AI platform
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.yourwebsite.com/support
// @contact.email  support@yourwebsite.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	server := server.SetupServer()
	cfg := config.Load()
	db, err := config.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	server.LoadHTMLGlob("templates/*")
	server.Static("/static", "./static")

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434"
	}

	ragRepo := repositories.NewRagRepository(db)
	ragService := services.NewRagService(ragRepo)
	ragHandler := handlers.NewRagHandler(ragService)

	chatRepo := repositories.NewChatRepository(ollamaURL,	cfg.OllamaDefaultModel, cfg.OpenAIAPIKey, cfg.GeminiAPIKey)
	searchRepo := utils.NewWebSearchRepository()
	chatService := services.NewChatService(chatRepo, searchRepo, ragRepo)
	chatHandler := handlers.NewChatHandler(chatService, db)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	deps := &routes.Dependencies{
		ChatHandler: chatHandler,
		UserHandler: userHandler,
		RagHandler: ragHandler,
	}

	routes.SetupRoutes(server, deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("server starting on port %s", port)
	log.Println("Swagger documentation available at: http://localhost:8080/swagger/index.html")
	server.Run("0.0.0.0:" + port)
}

