package utils

import(
	"net/http"
	"github.com/gin-gonic/gin"
)


func SwaggerRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{

		"GET https://localhost/api/documents": "Busca todos os dados do banco",
		"GET https://localhost:/api/documents/search?q= &limit=3":	"Procura no banco por similaridade. Colocar busca depois de q= 'sua busca'",
		"POST https://localhost/api/documents": "Envia texto para bando de dados",
		"POST https://localhost/api/chat": "Manda uma mensage para conversar com a IA",

		// "POST": "http://localhost:8080/api/post_code -->Envia texto para banco(Recomendado para o envio de codigo)",
		// "GETT.": "http://localhost:8080/api/get_code?q= &limit=3 --> Procura no banco por similaridade. Colocar busca depois de q= 'sua buscai",
		// "GETT_ALL":	"http://localhost:8080/get_all_code",
	})
}

