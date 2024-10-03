package main

import (
	"log"
	"post_service/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	r := gin.Default()
	routes.MappingRoute(r)
	r.Run(":3005")
}
