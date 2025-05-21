package main

import (
	"Assignment2_AdelKenesova/pkg/redis"
	"Assignment2_AdelKenesova/user_service/internal/db"
	"Assignment2_AdelKenesova/user_service/internal/service"
	pb "Assignment2_AdelKenesova/user_service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	redis.InitRedis()
	log.Println("Connected to Redis")

	db.InitDB()

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &service.UserService{})

	log.Println("UserService is running on port 50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
