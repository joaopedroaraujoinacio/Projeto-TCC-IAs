package server

import "github.com/gin-gonic/gin"


func SetupServer() *gin.Engine {
	server := gin.Default()
	server.GET("/test_server", func(c *gin.Context){
		c.String(200, "server working on port 8080")
	})
	return server
}
