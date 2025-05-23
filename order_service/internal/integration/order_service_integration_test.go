package integration

import (
	"Assignment2_AdelKenesova/order_service/internal/service"
	pb "Assignment2_AdelKenesova/order_service/proto"
	"context"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func startOrderServiceServer(t *testing.T) (*grpc.Server, pb.OrderServiceClient, func()) {
	lis, err := net.Listen("tcp", ":0") // любой свободный порт
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	orderSrv := &service.OrderService{}
	pb.RegisterOrderServiceServer(grpcServer, orderSrv)

	go grpcServer.Serve(lis)

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	assert.NoError(t, err)

	client := pb.NewOrderServiceClient(conn)

	cleanup := func() {
		grpcServer.Stop()
		conn.Close()
		lis.Close()
	}

	return grpcServer, client, cleanup
}

func TestOrderIntegration_CreateAndGet(t *testing.T) {
	_, client, cleanup := startOrderServiceServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Создание заказа
	orderResp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId: 1,
		Items: []*pb.OrderItem{
			{ProductId: 101, Quantity: 2, Price: 10.0},
			{ProductId: 102, Quantity: 1, Price: 20.0},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, orderResp)
	assert.Equal(t, float64(40.0), orderResp.Order.TotalPrice)

	// Получение заказа по ID
	getResp, err := client.GetOrder(ctx, &pb.GetOrderRequest{Id: orderResp.Order.Id})
	assert.NoError(t, err)
	assert.Equal(t, orderResp.Order.Id, getResp.Order.Id)
	assert.Equal(t, orderResp.Order.TotalPrice, getResp.Order.TotalPrice)
}

func TestOrderIntegration_ListOrders(t *testing.T) {
	_, client, cleanup := startOrderServiceServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Создаем несколько заказов
	for i := 0; i < 2; i++ {
		_, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
			UserId: 100 + uint64(i),
			Items: []*pb.OrderItem{
				{ProductId: 1, Quantity: 1, Price: 10.0},
			},
		})
		assert.NoError(t, err)
	}

	// Получаем все заказы
	listResp, err := client.ListOrders(ctx, &pb.Empty{})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResp.Orders), 2)
}

func TestOrderIntegration_DeleteOrder(t *testing.T) {
	_, client, cleanup := startOrderServiceServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Сначала создаем заказ
	createResp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId: 500,
		Items: []*pb.OrderItem{
			{ProductId: 1, Quantity: 1, Price: 10.0},
		},
	})
	assert.NoError(t, err)

	// Удаляем заказ
	_, err = client.DeleteOrder(ctx, &pb.DeleteOrderRequest{
		Id: createResp.Order.Id,
	})
	assert.NoError(t, err)

	// Пытаемся получить удалённый заказ
	_, err = client.GetOrder(ctx, &pb.GetOrderRequest{Id: createResp.Order.Id})
	assert.Error(t, err)
}
