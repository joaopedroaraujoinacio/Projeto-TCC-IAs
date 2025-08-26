package handlers

import (
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang_crud/services"
	"golang_crud/models"
)

type TestHandler struct {
	service *services.TestService
}

func NewTestHandler(service *services.TestService) *TestHandler {
	return  &TestHandler{service: service}
}

func (h *TestHandler) GetTests(c *gin.Context) {
	tests, err := h.service.GetAllTests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, 
		gin.H{"error": "Failed to get tests"})
		return
	}
	
	c.JSON(http.StatusOK, tests)
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		if err := services.CreateUser(db, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user": user,



		})
	}
}
