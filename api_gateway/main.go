package main

import (
	"Assignment2_AdelKenesova/api_gateway/handlers"
	"google.golang.org/grpc"
	"log"
	"net/http"

	"Assignment2_AdelKenesova/api_gateway/routes"
)

func main() {
	// gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}
	defer conn.Close()

	handlers.InitInventoryClient(conn)

	orderConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to OrderService: %v", err)
	}
	defer orderConn.Close()

	handlers.InitOrderClient(orderConn)

	userConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to UserService: %v", err)
	}
	defer userConn.Close()

	handlers.InitUserClient(userConn)

	// Router
	router := routes.SetupRouter()

	log.Println("API Gateway running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
