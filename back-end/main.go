package main

import (
	"os"
	"log"
	// _"golang_crud/docs"
	"golang_crud/config"
	"golang_crud/routes"
	"golang_crud/database"
	"github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"nginx"})

	r.GET("/healthy", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status": "healthy",
			"https":  true,
			"proxy":  "nginx",
		})
	})

	// routes.Setup(router, db)
	routes.SetupRoutes(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server starting on port %s", port)
	r.Run("0.0.0.0:" + port)

}
