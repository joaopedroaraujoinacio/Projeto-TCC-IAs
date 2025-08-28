package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"golang_crud/config"
	"golang_crud/database"
	"golang_crud/routes"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	router := gin.Default()
	routes.Setup(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)

}
