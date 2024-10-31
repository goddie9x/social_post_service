package main

import (
	"log"
	"sync"

	"post_service/internal/grpc"
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
	httpPort := dotenv.GetEnvOrDefaultValue("API_PORT", "3005")

	grpcPort := dotenv.GetEnvOrDefaultValue("GRPC_PORT", "50051")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP server on port %s", httpPort)
		if err := r.Run(":" + httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Printf("Starting gRPC server on port %s", grpcPort)
		grpc.StartGRPCServer(grpcPort)
	}()

	wg.Wait()
}
