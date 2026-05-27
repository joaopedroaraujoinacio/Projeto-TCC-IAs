package handlers

import (
	"go-project/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("signup bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if err := h.svc.Signup(&user); err != nil {
		log.Printf("signup error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("login bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	loggedIn, err := h.svc.Login(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": gin.H{
			"id":    loggedIn.ID,
			"email": loggedIn.Email,
			"role":  loggedIn.Role,
		},
	})
}
