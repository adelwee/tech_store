package service

import (
	"Assignment2_AdelKenesova/order_service/internal/db"
	"Assignment2_AdelKenesova/order_service/internal/models"
	pb "Assignment2_AdelKenesova/order_service/proto"
	"Assignment2_AdelKenesova/pkg/events"
	"Assignment2_AdelKenesova/pkg/nats"
	rdb "Assignment2_AdelKenesova/pkg/redis"
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	var total float64
	var orderItems []models.OrderItem

	for _, item := range req.Items {
		total += item.Price * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ProductID: uint(item.ProductId),
			Quantity:  uint(item.Quantity),
			Price:     item.Price,
		})
	}

	order := models.Order{
		UserID:     uint(req.UserId),
		TotalPrice: total,
		Status:     "pending",
		OrderItems: orderItems,
	}

	if err := db.DB.Create(&order).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create order: %v", err)
	}

	// Публикация события order.created в NATS
	go func() {
		event := events.OrderCreatedEvent{
			OrderID: uint64(order.ID),
			Items:   []events.OrderItem{},
		}

		for _, item := range order.OrderItems {
			event.Items = append(event.Items, events.OrderItem{
				ProductID: uint64(item.ProductID),
				Quantity:  uint64(item.Quantity),
			})
		}

		if err := nats.Publish("order.created", event); err != nil {
			log.Printf(" Failed to publish order.created: %v", err)
		} else {
			log.Printf(" Published order.created for Order ID %d", order.ID)
		}
	}()

	var pbItems []*pb.OrderItem
	for _, item := range order.OrderItems {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId: uint64(item.ProductID),
			Quantity:  uint64(item.Quantity),
			Price:     item.Price,
		})
	}

	return &pb.OrderResponse{
		Order: &pb.Order{
			Id:         uint64(order.ID),
			UserId:     uint64(order.UserID),
			Items:      pbItems,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	var order models.Order

	if err := db.DB.Preload("OrderItems").First(&order, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Order not found")
	}

	var pbItems []*pb.OrderItem
	for _, item := range order.OrderItems {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId: uint64(item.ProductID),
			Quantity:  uint64(item.Quantity),
			Price:     item.Price,
		})
	}

	return &pb.OrderResponse{
		Order: &pb.Order{
			Id:         uint64(order.ID),
			UserId:     uint64(order.UserID),
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt.Format(time.RFC3339),
			Items:      pbItems,
		},
	}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.Empty, error) {
	if err := db.DB.Delete(&models.Order{}, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete order: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, _ *pb.Empty) (*pb.ListOrdersResponse, error) {
	cacheKey := "orders:all"
	client := rdb.GetClient()

	//  Попробовать получить из Redis
	cached, err := client.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedOrders []*pb.Order
		if err := json.Unmarshal([]byte(cached), &cachedOrders); err == nil {
			log.Println("Returned orders from Redis cache")
			return &pb.ListOrdersResponse{Orders: cachedOrders}, nil
		}
	}

	//  Получить из базы, если в кэше нет
	var orders []models.Order
	if err := db.DB.Preload("OrderItems").Find(&orders).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list orders")
	}

	var pbOrders []*pb.Order
	for _, o := range orders {
		var items []*pb.OrderItem
		for _, item := range o.OrderItems {
			items = append(items, &pb.OrderItem{
				ProductId: uint64(item.ProductID),
				Quantity:  uint64(item.Quantity),
				Price:     item.Price,
			})
		}
		pbOrders = append(pbOrders, &pb.Order{
			Id:         uint64(o.ID),
			UserId:     uint64(o.UserID),
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedAt:  o.CreatedAt.Format(time.RFC3339),
			Items:      items,
		})
	}

	//  Сохранить в Redis
	data, err := json.Marshal(pbOrders)
	if err == nil {
		err = client.Set(ctx, cacheKey, data, 0).Err()
		if err == nil {
			log.Println("Orders cached in Redis")
		} else {
			log.Println("Failed to cache orders in Redis:", err)
		}
	}

	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}
