package middleware

import (
	"strings"
	"net/http"
	"go-project/utils"
	"github.com/gin-gonic/gin"
)


func Authenticate(context *gin.Context) {

	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token not found"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	context.Set("userId", userId)
	context.Next()
}

