package utils

import(


	"net/http"
	"github.com/gin-gonic/gin"
)


func SwaggerRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
// 		http://localhost:8080/api/documents
// http://localhost:8080/api/post_code
// http://localhost:8080/api/documents/search?q= &limit=3
			"message": "Hello world",
			"POST": "http://localhost:8080/api/documents --> Envia texto para banco de dados", 
	})
}
