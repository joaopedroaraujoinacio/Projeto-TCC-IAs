package routes

import (
	"net/http"
	"database/sql"
	"golang_crud/handlers"
	"golang_crud/repositories"
	"golang_crud/services"
	"github.com/gin-gonic/gin"
)


func Setup(r *gin.Engine, db *sql.DB) {
	testRepo := repositories.NewTestRepository(db)
	testService := services.NewTestService(testRepo)
	testHandler := handlers.NewTestHandler(testService)

	documentHandlers := &handlers.DocumentHandlers{DB: db}

	api := r.Group("/api")
	{
		api.GET("/tests", testHandler.GetTests)
		api.POST("/users", handlers.CreateUser(db))
		api.GET("/users/:id", handlers.GetUser(db))
		api.POST("/documents", handlers.CreateDocument(db))
		api.GET("/documents/search", documentHandlers.SearchDocuments)
		api.GET("/documents", documentHandlers.GetAllDocuments)

		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello world",
			})
		})
	}
}

