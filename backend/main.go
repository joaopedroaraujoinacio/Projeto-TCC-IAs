package main

import (
	"os"
	"log"
	"go-project/config"
	"go-project/routes"
	"go-project/server"

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

	routes.SetupRoutes(server, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("server starting on port %s", port)
	log.Println("Swagger documentation available at: http://localhost:8080/swagger/index.html")
	server.Run("0.0.0.0:" + port)
}

