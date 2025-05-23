package integration

import (
	"Assignment2_AdelKenesova/inventory_service/internal/service"
	pb "Assignment2_AdelKenesova/inventory_service/proto"
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func startTestGRPCServer(t *testing.T) (*grpc.Server, pb.InventoryServiceClient, func()) {
	lis, err := net.Listen("tcp", ":0") // порт выберется автоматически
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, &service.InventoryService{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}

	client := pb.NewInventoryServiceClient(conn)

	cleanup := func() {
		grpcServer.Stop()
		conn.Close()
	}

	return grpcServer, client, cleanup
}

func TestInventoryIntegration_CreateAndGet(t *testing.T) {
	_, client, cleanup := startTestGRPCServer(t)
	defer cleanup()

	// Создание продукта
	createReq := &pb.CreateProductRequest{
		Name:        "Integration Product",
		Description: "Integration Description",
		Price:       99.99,
	}
	createRes, err := client.CreateProduct(context.Background(), createReq)
	assert.NoError(t, err)
	assert.Equal(t, "Integration Product", createRes.Product.Name)

	// Получение продукта по ID
	getRes, err := client.GetProduct(context.Background(), &pb.GetProductRequest{
		Id: createRes.Product.Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, createRes.Product.Id, getRes.Product.Id)
	assert.Equal(t, "Integration Product", getRes.Product.Name)
}
