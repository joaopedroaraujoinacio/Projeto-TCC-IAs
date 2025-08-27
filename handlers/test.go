package handlers

import (
	"net/http"
	"strconv"
	"database/sql"
	"golang_crud/models"
	"golang_crud/services"
	"github.com/gin-gonic/gin"
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

func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := services.GetUserByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
