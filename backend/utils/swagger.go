package utils

import(
	"net/http"
	"github.com/gin-gonic/gin"
)


func SwaggerRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"POST https://localhost/api/chat": "Manda uma mensage para conversar com a IA",
	})
}

		// "GET https://localhost/api/documents": "Busca todos os dados do banco",
		// "GET https://localhost:/api/documents/search?q= &limit=3":	"Procura no banco por similaridade. Colocar busca depois de q= 'sua busca'",
		// "POST https://localhost/api/documents": "Envia texto para bando de dados",

