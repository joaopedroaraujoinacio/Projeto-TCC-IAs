package main

import (
	"os"
	"log"
	"net/http"
	"go-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello World with HTMX",
		})
	})

	// HTMX endpoint that returns HTML fragment
	r.GET("/hello", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{
			"message": "Hello World from HTMX!",
		})
	})

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("server starting on port %s", port)
	r.Run("0.0.0.0:" + port)
}

