package main

import (
	"os"
	"log"
	"go-project/config"
	"go-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	cfg := config.Load()
	db, err := config.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()


	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "/static")

	routes.SetupRoutes(r)
	//r, db

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("server starting on port %s", port)
	r.Run("0.0.0.0:" + port)
}

