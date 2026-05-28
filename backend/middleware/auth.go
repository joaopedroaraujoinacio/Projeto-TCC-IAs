package middleware

import (
	"net/http"
	"go-project/utils"
	"strings"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized"})
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		userID, role, err := utils.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
