package main

import (
	"Assignment2_AdelKenesova/order_service/internal/db"
	"Assignment2_AdelKenesova/order_service/internal/service"
	pb "Assignment2_AdelKenesova/order_service/proto"
	nats "Assignment2_AdelKenesova/pkg/nats"
	"Assignment2_AdelKenesova/pkg/redis"

	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	redis.InitRedis()
	log.Println(" Connected to Redis")

	//  Подключение к NATS
	if err := nats.InitNATS("nats://localhost:4222"); err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	//  БД
	db.InitDB()
	db.Migrate()

	//  gRPC сервер
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &service.OrderService{})

	log.Println(" OrderService is running on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
