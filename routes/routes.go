package routes

import (
	"database/sql"
	"golang_crud/handlers"
	"golang_crud/repositories"
	"golang_crud/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, db *sql.DB) {
	testRepo := repositories.NewTestRepository(db)
	testService := services.NewTestService(testRepo)
	testHandler := handlers.NewTestHandler(testService)

	api := r.Group("/api/v1")
	{
		api.GET("/tests", testHandler.GetTests)

		api.POST("/users", handlers.CreateUser(db))

		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello world",
			})
		})

	}
}

