package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	PrepareRoutes(server)

	server.Run(":8080")
}
