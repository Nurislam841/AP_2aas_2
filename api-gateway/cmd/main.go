package main

import (
	"log"

	"api-gateway/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)

	log.Println("API Gateway running on :8080")
	r.Run(":8080")
}
