package main

import (
	"os"
	"log"
	"golang_crud/config"
	"golang_crud/routes"
	"golang_crud/database"
	"github.com/gin-gonic/gin"
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

	routes.SetupRoutes(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run("0.0.0.0:" + port)
}

