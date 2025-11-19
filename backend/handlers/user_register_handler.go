package handlers

import (
	"net/http"
	"database/sql"
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

// func login(context *gin.Context) {
// 	var user services.User
//
// 	err := context.ShouldBindJSON(&user)
//
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid credentials", "error": err.Error()} )
// 		return 
// 	}
//
// 	err = user.ValidateCredentials()
// 	if err != nil {
// 		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"} )
// 		return
// 	}
//
// 	token, err := utils.GenerateToken(user.Email, user.ID)
// 	if err != nil {
// 		log.Printf("token generation failed for user %d: %v", user.ID, err)
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not authenticate user", "error": err} )
// 		return 
// 	}
//
// 	context.JSON(http.StatusOK, gin.H{"message": "login successfull", "token": token })
// }
//
