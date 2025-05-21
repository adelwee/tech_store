package main

import (
	"Assignment2_AdelKenesova/inventory_service/internal/db"
	"Assignment2_AdelKenesova/inventory_service/internal/service"
	pb "Assignment2_AdelKenesova/inventory_service/proto"
	"Assignment2_AdelKenesova/pkg/nats"
	"Assignment2_AdelKenesova/pkg/redis"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// Подключение к Redis
	redis.InitRedis()
	log.Println(" Connected to Redis")

	// Подключение к NATS
	if err := nats.InitNATS("nats://localhost:4222"); err != nil {
		log.Fatalf(" Failed to connect to NATS: %v", err)
	}
	nc := nats.GetConn()

	// База данных
	db.InitDB()
	db.Migrate()

	// gRPC-сервис
	serviceImpl := &service.InventoryService{}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf(" Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, serviceImpl)

	// gRPC-клиент для NATS
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf(" Failed to connect to self gRPC: %v", err)
	}
	nats.SetInventoryClient(pb.NewInventoryServiceClient(conn))

	// Подписка на события
	if err := nats.SubscribeToProductCreated(nc); err != nil {
		log.Fatalf(" Failed to subscribe to product.created: %v", err)
	}
	if err := nats.SubscribeToOrderCreated(nc); err != nil {
		log.Fatalf(" Failed to subscribe to order.created: %v", err)
	}

	log.Println(" InventoryService is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(" Failed to serve: %v", err)
	}
}
