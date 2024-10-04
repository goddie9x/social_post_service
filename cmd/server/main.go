package main

import (
	"log"
	"post_service/internal/routes"
	"post_service/pkg/configs"
	"post_service/pkg/dotenv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	r := gin.Default()
	routes.MappingRoute(r)
	d := configs.DiscoveryServerConnect{}
	d.ConnectToEurekaDiscoveryServer()
	port := dotenv.GetEnvOrDefaultValue("API_PORT", "3005")
	r.Run(":" + port)
}
