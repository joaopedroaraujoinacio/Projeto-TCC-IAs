package utils

import(


	"net/http"
	"github.com/gin-gonic/gin"
)


func SwaggerRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
			"GET":	"http://localhost:8080/api/documents --> Busca todos os dados do banco",
			"POST": "http://localhost:8080/api/post_code -->Envia texto para banco(Recomendado para o envio de codigo)",
			"POST.": "http://localhost:8080/api/documents/search?q= &limit=3 --> Procura no banco por similaridade. Colocar busca depois de q= 'sua buscai",
			"POST..": "http://localhost:8080/api/documents --> Envia texto para banco de dados", 
			"message":	"Hello world",

			"GETT.": "http://localhost:8080/api/get_code?q= &limit=3 --> Procura no banco por similaridade. Colocar busca depois de q= 'sua buscai",
			//fix
	})
}
