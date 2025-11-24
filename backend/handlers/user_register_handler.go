package handlers

import (
	"log"
	"net/http"
	"database/sql"
	"go-project/utils"
	"go-project/services"
	"github.com/gin-gonic/gin"
)


func Signup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	var user services.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"} )
		return
	}

	err = user.Save(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error"} )
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	var user services.User

	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid credentials", "error": err.Error()} )
		return 
	}

	err = user.ValidateCredentials(db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"} )
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		log.Printf("token generation failed for user %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not authenticate user", "error": err} )
		return 
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successfull", "token": token })
	}
}

